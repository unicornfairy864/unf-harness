package client

import (
	"net/http"
)

func Init(client *LLMClient) {
	client.httpClient = &http.Client{Timeout: connectionTimeout}
}
