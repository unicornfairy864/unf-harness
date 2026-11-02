package tools

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"unicode/utf8"

	"github.com/unicornfairy864/unf-harness/constants"
	"github.com/unicornfairy864/unf-harness/model"
	"github.com/unicornfairy864/unf-harness/utils"
)

type toolWriteFile struct {
	Function model.Function
	Params   paramsWriteFile
}

func (t *toolWriteFile) Name() string {
	return t.Function.Name
}

func (t *toolWriteFile) GetTool() model.Function {
	return t.Function
}

func (t *toolWriteFile) ReceiveParam(message string) error {
	if strings.TrimSpace(message) == "" {
		message = "{}"
	}
	params := paramsWriteFile{}
	err := json.Unmarshal([]byte(message), &params)
	if err != nil {
		return err
	}
	t.Params = params
	return nil
}

type paramsWriteFile struct {
	Mode     string `json:"mode"`
	FilePath string `json:"file_path"`
	Content  string `json:"content"`
}

const (
	modeOverwrite = "overwrite"
	modeAppend    = "append"

	maxWriteLength = 81920
	writeFilePerm  = 0o644
	writeDirPerm   = 0o755
)

const maxConfirmRetry = 2

func (t *toolWriteFile) Execute() (string, error) {
	filePath := strings.TrimSpace(t.Params.FilePath)
	if filePath == "" {
		return "", errors.New("file_path is required")
	}

	content := t.Params.Content
	if content == "" {
		return "", errors.New("content is empty")
	}
	if n := utf8.RuneCountInString(content); n > maxWriteLength {
		return "", fmt.Errorf("content too long: %d chars, max %d", n, maxWriteLength)
	}

	// mode 归一化：空串当 overwrite。只认这两个值，写错就报错，不猜
	mode := strings.TrimSpace(t.Params.Mode)
	flags := os.O_WRONLY | os.O_CREATE
	switch mode {
	case "", modeOverwrite:
		mode = modeOverwrite
		flags |= os.O_TRUNC
	case modeAppend:
		flags |= os.O_APPEND
	default:
		return "", fmt.Errorf("unknown mode: %q, want %q or %q", mode, modeOverwrite, modeAppend)
	}

	// 相对路径基准 = config 里的 workplace，跟 read_file 同一套
	path := resolvePath(filePath)

	// —— 写之前先摆给人看 ——
	// 第一行 mode + 路径；第二行内容预览。
	// 预览里的换行转成字面 \n，否则多行内容会把"第二行"撑成几十行
	fmt.Printf("[write_file] %s %s\n", mode, path)
	fmt.Println(utils.Truncate(strings.ReplaceAll(content, "\n", `\n`), constants.TruncateMaxLength))

	if !confirmWrite(content) {
		return "", errors.New("write rejected by user")
	}

	// 父目录不存在就建出来：模型很爱写 static/example.txt 这种
	if dir := filepath.Dir(path); dir != "" {
		if err := os.MkdirAll(dir, writeDirPerm); err != nil {
			return "", fmt.Errorf("create directory failed: %s: %v", dir, err)
		}
	}

	if info, err := os.Stat(path); err == nil && info.IsDir() {
		return "", fmt.Errorf("%s is a directory, not a file", filePath)
	}

	f, err := os.OpenFile(path, flags, writeFilePerm)
	if err != nil {
		switch {
		case os.IsPermission(err):
			return "", fmt.Errorf("permission denied: %s", filePath)
		default:
			return "", fmt.Errorf("open file failed: %s: %v", filePath, err)
		}
	}

	if _, err := f.WriteString(content); err != nil {
		_ = f.Close()
		return "", fmt.Errorf("write file failed: %s: %v", filePath, err)
	}
	// 写文件的 Close 不能 defer 掉错误：缓冲区可能还没落盘
	if err := f.Close(); err != nil {
		return "", fmt.Errorf("close file failed: %s: %v", filePath, err)
	}

	return fmt.Sprintf("%s ok: %s (%d chars)", mode, filePath, utf8.RuneCountInString(content)), nil
}

// confirmWrite 写入前问一次人：y 允许、n 拒绝、d 看完整内容再问一遍。
// 非 y/n/d 的输入（含直接回车）累计 maxConfirmRetry 次 → 自动拒绝。
func confirmWrite(content string) bool {
	invalid := 0
	for {
		fmt.Print("[write_file] allow write? (y=yes / n=no / d=detail) > ")

		line, ok := utils.InputLine()
		if !ok { // EOF（Ctrl+Z / Ctrl+D）当拒绝处理，别让 Agent 卡在这儿
			fmt.Println()
			return false
		}

		switch strings.ToLower(strings.TrimSpace(line)) {
		case "y":
			return true
		case "n":
			return false
		case "d":
			fmt.Println(content) // 完整内容，不截断
		default:
			invalid++
			if invalid >= maxConfirmRetry {
				fmt.Println("[write_file] too many invalid inputs, rejected.")
				return false
			}
			fmt.Printf("[write_file] please answer y / n / d (%d/%d)\n", invalid, maxConfirmRetry)
		}
	}
}

var ToolWriteFile = toolWriteFile{
	Function: model.Function{
		Name:        "write_file",
		Description: "写入工具。把 content 写入 file_path 指定的文件。mode=overwrite 覆盖写入（文件不存在则创建，已存在则清空后写入），mode=append 追加到文件末尾；单次最多 8192 字符，更长的内容分多次 append。",
		Strict:      false,
		Parameters: model.Schema{
			Type: "object",
			Properties: map[string]model.Property{
				"mode": {
					Type:        "string",
					Description: "编辑模式。overwrite：覆盖；append：追加到末尾。默认为 overwrite",
					Enum:        []string{modeOverwrite, modeAppend},
				},
				"file_path": {
					Type:        "string",
					Description: "文件路径。若以\"大写字母:\"开头则为绝对路径，否则为相对路径",
				},
				"content": {
					Type:        "string",
					Description: "要写入的内容主体，最多 81920 字符",
				},
			},
			Required: []string{"file_path", "content"},
		},
	},
}
