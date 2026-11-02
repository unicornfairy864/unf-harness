package client

import (
	"net/http"

	"github.com/unicornfairy864/unf-harness/config"
	"github.com/unicornfairy864/unf-harness/constants"
	"github.com/unicornfairy864/unf-harness/model"
	"github.com/unicornfairy864/unf-harness/tool"
)

func Init(client *LLMClient) {
	client.httpClient = &http.Client{Timeout: config.UH_Config.LLM.Timeout}
}

func InitRequestBody(enableTool bool) *model.RequestBody {
	req := &model.RequestBody{
		Model:           config.UH_Config.LLM.Model,
		ReasoningEffort: config.UH_Config.LLM.ReasoningEffort,
		MaxTokens:       config.UH_Config.LLM.MaxTokens,
		Temperature:     config.UH_Config.LLM.Temperature,
	}
	req.SetResponseFormatText()
	if config.UH_Config.LLM.Thinking {
		req.EnableThinking()
	} else {
		req.DisableThinking()
	}

	// Message
	req.AppendText(constants.RoleSystem, config.UH_Config.LLM.SystemPrompt)

	// Tool
	if enableTool {
		for _, toolName := range config.UH_Config.Tools.Enabled {
			t, ok := tool.Tools[toolName]
			if !ok {
				continue
			}
			req.RegisterTool(model.ToolToFuncTool(t.GetTool()))
		}
	}

	return req
}
