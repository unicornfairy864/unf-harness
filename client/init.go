package client

import (
	"net/http"

	"github.com/unicornfairy864/unf-harness/config"
)

func Init(client *LLMClient) {
	client.httpClient = &http.Client{Timeout: config.UH_Config.LLM.Timeout}
}
