package client

import (
	"net/http"
	"time"
)

var UH_Client LLMClient

const (
	connectionTimeout = 30 * time.Second
)

type LLMClient struct {
	httpClient *http.Client
}
