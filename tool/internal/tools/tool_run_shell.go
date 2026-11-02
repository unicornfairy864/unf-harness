package tools

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/unicornfairy864/unf-harness/config"
	"github.com/unicornfairy864/unf-harness/model"
	"github.com/unicornfairy864/unf-harness/utils"
)

type toolRunShell struct {
	Function model.Function
	Params   paramsRunShell
}

func (t *toolRunShell) Name() string { return t.Function.Name }

func (t *toolRunShell) GetTool() model.Function { return t.Function }

func (t *toolRunShell) ReceiveParam(message string) error {
	if strings.TrimSpace(message) == "" {
		message = "{}"
	}
	params := paramsRunShell{}
	if err := json.Unmarshal([]byte(message), &params); err != nil {
		return err
	}
	t.Params = params
	return nil
}

type paramsRunShell struct {
	Command     string `json:"command"`
	Description string `json:"description"`
}

const (
	maxCommandLength           = 8192
	defaultShellTimeout        = 60 * time.Second
	defaultShellMaxOutputBytes = 16384
	shellWaitDelay             = 5 * time.Second
)

func (t *toolRunShell) Execute() (string, error) {
	command := strings.TrimSpace(t.Params.Command)
	if command == "" {
		return "", errors.New("command is required")
	}
	if len(command) > maxCommandLength {
		return "", fmt.Errorf("command too long: %d bytes, max %d", len(command), maxCommandLength)
	}

	// 目录先验一遍：不存在就让 cmd.Run 去报，错误会很难看
	dir := workplaceRoot()
	if info, err := os.Stat(dir); err != nil || !info.IsDir() {
		return "", fmt.Errorf("workplace %q is not a usable directory", dir)
	}

	risk := scanRisk(command)

	// —— 摆给人看：目录、意图、完整命令 ——
	fmt.Printf("[run_shell] cwd: %s\n", dir)
	if intent := strings.TrimSpace(t.Params.Description); intent != "" {
		fmt.Printf("[run_shell] intent: %s\n", intent)
	}
	fmt.Printf("[run_shell] $ %s\n", command)
	if risk != "" {
		fmt.Printf("[run_shell] WARNING: %s\n", risk)
	}

	if !confirmShell(risk != "") {
		return "", errors.New("shell command rejected by user")
	}

	return runShell(command, dir)
}

// runShell 在 dir 下跑一条命令。超时、输出上限、进程树清理都在这儿兜。
func runShell(command, dir string) (string, error) {
	timeout := config.UH_Config.Tools.RunShell.Timeout
	if timeout <= 0 {
		timeout = defaultShellTimeout
	}
	maxOut := config.UH_Config.Tools.RunShell.MaxOutputBytes
	if maxOut <= 0 {
		maxOut = defaultShellMaxOutputBytes
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.CommandContext(ctx, "cmd.exe", "/c", command)
	} else {
		cmd = exec.CommandContext(ctx, "/bin/sh", "-c", command)
	}
	cmd.Dir = dir

	// 不接 stdin：命令要交互就直接拿到 EOF。
	// 不这么做的话，子进程会跟 CLI 抢同一个 stdin，
	// 后面那个确认框就会读到被抢剩下的碎片。
	cmd.Stdin = nil

	// 超时时默认只 kill 主进程，cmd.exe 底下的孙进程会变孤儿。
	// Windows 用 taskkill /T 连整棵树一起收；其他平台 kill 主进程就够。
	cmd.Cancel = func() error {
		if runtime.GOOS != "windows" {
			return cmd.Process.Kill()
		}
		pid := strconv.Itoa(cmd.Process.Pid)
		if err := exec.Command("taskkill", "/F", "/T", "/PID", pid).Run(); err != nil {
			return cmd.Process.Kill() // taskkill 不在就退回单杀
		}
		return nil
	}
	// 进程退了但输出管道没关时，Wait 最多再等这么久就强行收尾
	cmd.WaitDelay = shellWaitDelay

	// stdout 和 stderr 指向同一个 writer：Go 会开两条管道、
	// 两个 goroutine 并发写，所以 limitedWriter 里面必须上锁
	out := &limitedWriter{limit: maxOut}
	cmd.Stdout = out
	cmd.Stderr = out

	err := cmd.Run()

	var sb strings.Builder
	switch {
	case errors.Is(ctx.Err(), context.DeadlineExceeded):
		fmt.Fprintf(&sb, "TIMEOUT after %s (process tree killed)\n", timeout)
	case err != nil:
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			// 非零退出码不是工具错误：编译报错、grep 没命中都是正常结果，
			// 输出得原样给模型看，所以照常返回内容
			fmt.Fprintf(&sb, "exit code: %d\n", exitErr.ExitCode())
		} else {
			return "", fmt.Errorf("run command failed: %v", err)
		}
	default:
		sb.WriteString("exit code: 0\n")
	}

	sb.WriteString("--- output ---\n")
	sb.WriteString(out.String())
	if out.Truncated() {
		fmt.Fprintf(&sb, "\n...[truncated: kept %d of %d bytes]", maxOut, out.Total())
	}
	return sb.String(), nil
}

// limitedWriter 只保留前 limit 字节，多余的丢掉但照常计数。
// 命令疯狂输出时内存不能跟着一起涨。
// Write 永远返回 len(p), nil——返回短写会让子进程直接吃到 write 错误然后崩。
type limitedWriter struct {
	mu    sync.Mutex
	buf   bytes.Buffer
	limit int
	total int
}

func (w *limitedWriter) Write(p []byte) (int, error) {
	w.mu.Lock()
	defer w.mu.Unlock()

	w.total += len(p)
	if room := w.limit - w.buf.Len(); room > 0 {
		if len(p) <= room {
			_, _ = w.buf.Write(p)
		} else {
			_, _ = w.buf.Write(p[:room])
		}
	}
	return len(p), nil
}

// String 返回保留下的输出。截断点可能切碎一个多字节字符，顺手修掉。
func (w *limitedWriter) String() string {
	w.mu.Lock()
	defer w.mu.Unlock()
	return strings.ToValidUTF8(w.buf.String(), "")
}

func (w *limitedWriter) Total() int {
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.total
}

func (w *limitedWriter) Truncated() bool {
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.total > w.buf.Len()
}

// confirmShell 执行前问一次人：y 允许、n 拒绝。
// dangerous=true 时不认单个 y，必须打完整 "yes"。
// 非法输入（含直接回车）累计 maxConfirmRetry 次 → 自动拒绝。
// maxConfirmRetry 在 tool_write_file.go 里，同包共用。
func confirmShell(dangerous bool) bool {
	invalid := 0
	for {
		if dangerous {
			fmt.Print(`[run_shell] DANGEROUS command. type "yes" to run, n to reject > `)
		} else {
			fmt.Print("[run_shell] run it? (y=yes / n=no) > ")
		}

		line, ok := utils.InputLine()
		if !ok { // EOF 当拒绝，别让 Agent 卡死在这儿
			fmt.Println()
			return false
		}

		answer := strings.ToLower(strings.TrimSpace(line))
		switch {
		case answer == "n":
			return false
		case dangerous && answer == "yes", !dangerous && answer == "y":
			return true
		}

		invalid++
		if invalid >= maxConfirmRetry {
			fmt.Println("[run_shell] too many invalid inputs, rejected.")
			return false
		}
		fmt.Printf("[run_shell] invalid input (%d/%d)\n", invalid, maxConfirmRetry)
	}
}

// riskyPatterns：命中就把命令标成高危，确认时要人打完整 "yes"。
// 只覆盖「不可逆 / 破坏性 / 提权 / 远程代码执行」四类，不封杀——
// 合法的清理、部署也要用它们，只是不让人顺手一个 y 就放过去。
// 子串匹配会误伤（grep "rm -rf" README.md 这种），代价是多打三个字母，能接受。
var riskyPatterns = []struct{ sub, why string }{
	{"rm -rf", "递归强删"},
	{"rm -fr", "递归强删"},
	{"rm -r", "递归删目录"},
	{"del /f", "强制删除"},
	{"del /s", "递归删除"},
	{"del /q", "静默删除"},
	{"rd /s", "递归删目录"},
	{"rmdir /s", "递归删目录"},
	{"format ", "格式化磁盘"},
	{"diskpart", "磁盘分区操作"},
	{"fdisk", "磁盘分区操作"},
	{"mkfs", "格式化文件系统"},
	{"dd if=", "裸设备写入"},
	{"shutdown", "关机 / 重启"},
	{"reg delete", "删注册表项"},
	{"bcdedit", "改启动配置"},
	{"cipher /w", "擦除磁盘空闲区"},
	{"takeown", "夺取文件所有权"},
	{"icacls", "改文件 ACL"},
	{"sudo ", "提权执行"},
	{"chmod 777", "放开全部权限"},
	{"chown ", "改文件属主"},
	{"| sh", "管道执行远程脚本"},
	{"| bash", "管道执行远程脚本"},
	{"| iex", "管道执行远程脚本"},
	{"| powershell", "管道执行 PowerShell"},
	{"-enc ", "编码执行的 PowerShell"},
	{"iex(", "执行动态字符串"},
	{":(){", "fork 炸弹"},
	{"%0|%0", "fork 炸弹"},
	{"taskkill /f", "强杀进程"},
	{"kill -9", "强杀进程"},
}

// scanRisk 返回命中的风险说明，没命中返回空串
func scanRisk(command string) string {
	lower := strings.ToLower(command)
	for _, r := range riskyPatterns {
		if strings.Contains(lower, r.sub) {
			return r.why
		}
	}
	return ""
}

var ToolRunShell = toolRunShell{
	Function: model.Function{
		Name: "run_shell",
		Description: "在 workplace 工作目录下执行一条 shell 命令（Windows 走 cmd.exe /c，其他平台走 /bin/sh -c）。" +
			"执行前会把完整命令交给用户确认，用户同意才真正运行，所以 description 要写清意图。" +
			"适合查看目录、跑构建和测试、git 查询等；" +
			"不适合需要交互输入或长期挂起的命令——stdin 是空的，超时会连进程树一起杀掉。" +
			"默认 60 秒超时，输出超出上限会被截断。",
		Strict: false,
		Parameters: model.Schema{
			Type: "object",
			Properties: map[string]model.Property{
				"command": {
					Type:        "string",
					Description: "要执行的完整命令，单条，最长 8192 字节",
				},
				"description": {
					Type:        "string",
					Description: "这条命令要做什么、为什么执行。会显示在用户确认框里，写清楚能少一次来回",
				},
			},
			Required: []string{"command", "description"},
		},
	},
}
