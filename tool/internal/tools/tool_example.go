package tools

import (
	"encoding/json"
	"errors"

	"github.com/unicornfairy864/unf-harness/model"
)

type toolExample struct {
	Function model.Function
	Params   paramsExample
}

func (t *toolExample) Name() string {
	return t.Function.Name
}

func (t *toolExample) GetTool() model.Function {
	return t.Function
}

func (t *toolExample) ReceiveParam(message string) error {
	params := paramsExample{}
	err := json.Unmarshal([]byte(message), &params)
	if err != nil {
		return err
	}
	t.Params = params
	return nil
}

type paramsExample struct {
	Content string `json:"content"`
}

func (t *toolExample) Execute() (string, error) {
	if t.Params.Content == "" {
		return "", errors.New("no content provided")
	}
	return t.Params.Content, nil
}

var ToolExample = toolExample{
	Function: model.Function{
		Name:        "example",
		Description: "示例工具。通过给定 content 固定返回原来的 content",
		Strict:      false,
		Parameters: model.Schema{
			Type: "object",
			Properties: map[string]model.Property{
				"content": {
					Type:        "string",
					Description: "需要原样返回的内容",
				},
			},
			Required: []string{"content"},
		},
	},
}
