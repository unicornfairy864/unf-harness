package main

import (
	"errors"
	"fmt"
	"io"

	"github.com/unicornfairy864/unf-harness/client"
	"github.com/unicornfairy864/unf-harness/config"
	"github.com/unicornfairy864/unf-harness/constants"
	"github.com/unicornfairy864/unf-harness/tool"
)

func main() {
	config.Init(&config.UH_Config)

	tool.Init()

	client.Init(&client.UH_Client)

	// AgentLoop

	var userPrompt string

	req := client.InitRequestBody()
	userPrompt = InputUserPrompt()
	for userPrompt != constants.CliExitCommand {
		req.AppendText("user", userPrompt)

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
		if data.Choices[0].FinishReason == constants.StopReasonToolCalls {
			fmt.Print(data.Choices[0].Message)
		} else {
			fmt.Println(data.Choices[0].Message)
			userPrompt = InputUserPrompt()
		}
	}
}

func InputUserPrompt() string {
	var str = ""
	fmt.Print(constants.CliInputPrompt)
	for {
		_, err := fmt.Scanf("%s\n", &str)
		if err != nil && errors.Is(err, io.EOF) {
			return constants.CliExitCommand
		}
		if err != nil || str == "" {
			fmt.Println(constants.CliInputAgainPrompt)
			fmt.Print(constants.CliInputPrompt)
			continue
		}
		break
	}
	return str
}
