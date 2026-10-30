package main

import (
	"fmt"

	"github.com/unicornfairy864/unf-harness/client"
	"github.com/unicornfairy864/unf-harness/config"
)

func main() {
	config.Init(&config.UH_Config)

	client.Init(&client.UH_Client)

	data, err := client.UH_Client.Request(
		config.UH_Config.LLM.Model,
		config.UH_Config.LLM.ReasoningEffort,
		config.UH_Config.LLM.SystemPrompt,
		config.UH_Config.LLM.UserPrompt,
		config.UH_Config.LLM.Thinking,
		false,
		config.UH_Config.LLM.MaxTokens)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println(*data)
	}
}
