package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/unicornfairy864/unf-harness/client"
	"github.com/unicornfairy864/unf-harness/config"
	"github.com/unicornfairy864/unf-harness/constants"
	"github.com/unicornfairy864/unf-harness/model"
	"github.com/unicornfairy864/unf-harness/tool"
	"github.com/unicornfairy864/unf-harness/utils"
)

func main() {
	config.Init(&config.UH_Config)

	tool.Init()

	client.Init(&client.UH_Client)

	// AgentLoop

	var userPrompt string

	req := client.InitRequestBody()
	userPrompt = InputUserPrompt()
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
			userPrompt = InputUserPrompt()
			req.AppendText(constants.RoleUser, userPrompt)
		}
	}
	fmt.Println(constants.CliExitPrompt)
}

var stdin = bufio.NewReader(os.Stdin)

// InputUserPrompt 读用户输入，支持 "\" 续行；EOF 返回 constants.CliExitCommand。
func InputUserPrompt() string {

	var sb strings.Builder

	for {
		if sb.Len() == 0 {
			fmt.Print(constants.CliUserPrompt)
		} else {
			fmt.Print(constants.CliContPrompt)
		}

		line, err := stdin.ReadString('\n')

		// 不用管是哪个平台：右侧的 \r 和 \n 一起删干净（cutset，不是子串）
		line = strings.TrimRight(line, "\r\n")

		// 续行判定：删掉行末空白后以 \ 结尾，且 \ 前面是空白（或整行就一个 \）
		body := strings.TrimRight(line, " \t")
		head := strings.TrimSuffix(body, `\`)
		isCont := strings.HasSuffix(body, `\`) &&
			(head == "" || strings.HasSuffix(head, " ") || strings.HasSuffix(head, "\t"))

		if isCont {
			sb.WriteString(strings.TrimRight(head, " \t")) // 续行符连同它前面的空白一起收走
			sb.WriteString("\n")
			if err != nil {
				break // 续行符后直接 EOF：攒下的内容照常返回
			}
			continue
		}

		// 非续行：纯空输入才提示重来（有累积内容时，空行算有效的一行）
		if sb.Len() == 0 && strings.TrimSpace(line) == "" {
			if err != nil {
				return constants.CliExitCommand
			}
			fmt.Println(constants.CliInputAgainPrompt)
			continue
		}

		sb.WriteString(line)
		break
	}

	return sb.String()
}
