package main

import (
	"encoding/json"
	"fmt"

	"github.com/unicornfairy864/unf-harness/client"
	"github.com/unicornfairy864/unf-harness/config"
	"github.com/unicornfairy864/unf-harness/constants"
	"github.com/unicornfairy864/unf-harness/model"
	"github.com/unicornfairy864/unf-harness/tool"
	"github.com/unicornfairy864/unf-harness/utils"
)

func main() {
	config.Init(&config.UH_Config)

	tool.Register()

	client.Init(&client.UH_Client)

	// AgentLoop

	var userPrompt string

	req := client.InitRequestBody(true)
	userPrompt = utils.InputUserPrompt()
	req.AppendText(constants.RoleUser, userPrompt)
	for userPrompt != constants.CliExitCommand {

		data, err := client.UH_Client.Request(req)
		if err != nil {
			fmt.Println(err.Error())
		}
		if data == nil {
			panic("Unknow Error")
		}
		req.AppendRawMessage(data.Choices[0].Message)

		if prompt, ok := constants.PromptStop[data.Choices[0].FinishReason]; ok {
			fmt.Println(prompt)
		} else {
			fmt.Println(data.Choices[0].FinishReason)
		}
		var messageResp model.MessageResponse
		err = json.Unmarshal(data.Choices[0].Message, &messageResp)
		if err != nil {
			panic(err.Error())
		}
		if messageResp.ReasoningContent != "" {
			fmt.Println(constants.CliThinkingPrompt + messageResp.ReasoningContent)
		}
		if messageResp.Content != "" {
			fmt.Println(constants.CliAssistantPrompt + messageResp.Content)
		}
		if data.Choices[0].FinishReason == constants.StopReasonToolCalls {
			for _, toolCall := range messageResp.ToolCalls {
				fmt.Printf("[ToolCall:%d-%s]Params:%s\n", toolCall.Index, toolCall.Function.Name, utils.Truncate(string(toolCall.Function.Arguments), constants.TruncateMaxLength))
				toolResp := "Tool:" + toolCall.Function.Name + " is undefined."
				if tool.Tools[toolCall.Function.Name] != nil {
					err = tool.Tools[toolCall.Function.Name].ReceiveParam(toolCall.Function.Arguments)
					if err != nil {
						toolResp = err.Error()
					} else {
						toolResp = tool.Execute(tool.Tools[toolCall.Function.Name])
					}
				}
				fmt.Printf("resp>" + utils.Truncate(toolResp, constants.TruncateMaxLength) + "\n")
				req.AppendToolCallWithID(toolCall.Id, toolResp)
			}
		} else {
			userPrompt = utils.InputUserPrompt()
			req.AppendText(constants.RoleUser, userPrompt)
		}
	}
	fmt.Println(constants.CliExitPrompt)
}
