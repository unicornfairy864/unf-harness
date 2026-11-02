package tools

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/unicornfairy864/unf-harness/model"
)

type toolReadFile struct {
	Function model.Function
	Params   paramsReadFile
}

func (t *toolReadFile) Name() string {
	return t.Function.Name
}

func (t *toolReadFile) GetTool() model.Function {
	return t.Function
}

func (t *toolReadFile) ReceiveParam(message string) error {
	if strings.TrimSpace(message) == "" {
		message = "{}"
	}
	params := paramsReadFile{}
	err := json.Unmarshal([]byte(message), &params)
	if err != nil {
		return err
	}
	t.Params = params
	return nil
}

type paramsReadFile struct {
	FilePath string `json:"file_path"`
	Offset   int64  `json:"offset"`
	Length   int64  `json:"length"`
}

var ToolReadFile = toolReadFile{
	Function: model.Function{
		Name:        "read_file",
		Description: "根据路径读取文件内容",
		Strict:      false,
		Parameters: model.Schema{
			Type: "object",
			Properties: map[string]model.Property{
				"file_path": {
					Type:        "string",
					Description: "需要读取的文件路径。若以\"大写字母:\"开头则为绝对路径，否则为相对路径",
				},
				"offset": {
					Type:        "integer",
					Description: "跳过的字符数，即从offset+1的字符起读。默认为 0 ",
				},
				"length": {
					Type:        "integer",
					Description: "最多读取的字符数。默认为 1024 ，最大 2048",
				},
			},
			Required: []string{"file_path"},
		},
	},
}

func (t *toolReadFile) Execute() (string, error) {
	filePath := strings.TrimSpace(t.Params.FilePath)
	if filePath == "" {
		return "", errors.New("file_path is required")
	}

	// 参数归一化：负 offset 当 0；length 非法给默认值，超上限直接截断（不报错）
	offset := t.Params.Offset
	if offset < 0 {
		offset = 0
	}
	length := t.Params.Length
	if length <= 0 {
		length = defaultReadLength
	}
	if length > maxReadLength {
		length = maxReadLength
	}

	path := resolvePath(filePath)

	f, err := os.Open(path)
	if err != nil {
		// 错误文本会经 executor 回灌给模型，所以分类写清楚，别甩 Windows 原始报错
		switch {
		case os.IsNotExist(err):
			return "", fmt.Errorf("file not found: %s", filePath)
		case os.IsPermission(err):
			return "", fmt.Errorf("permission denied: %s", filePath)
		default:
			return "", fmt.Errorf("open file failed: %s: %v", filePath, err)
		}
	}
	defer f.Close()

	// Windows 上 os.Open 能打开目录，必须自己拦
	info, err := f.Stat()
	if err != nil {
		return "", fmt.Errorf("stat file failed: %v", err)
	}
	if info.IsDir() {
		return "", fmt.Errorf("%s is a directory, not a file", filePath)
	}

	// 逐 rune 读，不用 Seek：按「字符」定位就没法用字节偏移跳
	reader := bufio.NewReaderSize(f, 64<<10)
	var sb strings.Builder
	sb.Grow(int(length) * 3) // 中文按 3 字节预留，省一次扩容

	var skipped, taken int64
	for taken < length {
		r, _, err := reader.ReadRune()
		if err == io.EOF {
			break // 读到尾不是错误，返回已有的就行
		}
		if err != nil {
			return "", fmt.Errorf("read file failed: %v", err)
		}
		if skipped < offset {
			skipped++
			continue
		}
		sb.WriteRune(r)
		taken++
	}

	out := sb.String()
	if offset == 0 {
		out = strings.TrimPrefix(out, "\uFEFF") // 顺手吃掉 UTF-8 BOM
	}
	return out, nil
}

const (
	defaultReadLength int64 = 1024
	maxReadLength     int64 = 2048
)

// isDriveAbsPath：第一个字符是大写字母、第二个是冒号 → 当绝对路径，其余一律当相对路径
func isDriveAbsPath(p string) bool {
	if len(p) < 2 || p[1] != ':' {
		return false
	}
	return p[0] >= 'A' && p[0] <= 'Z'
}
