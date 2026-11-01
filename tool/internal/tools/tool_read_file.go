package tools

import (
	"encoding/json"
	"errors"

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
	params := paramsReadFile{}
	err := json.Unmarshal([]byte(message), &params)
	if err != nil {
		return err
	}
	t.Params = params
	return nil
}

func (t *toolReadFile) Execute() (string, error) {
	return "", errors.New("tool not implemented")
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
