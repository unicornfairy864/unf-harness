package models

import "github.com/unicornfairy864/unf-harness/client/internal/constants"

type ChatRequestBody struct {
	Messages        []message `json:"messages"`
	Model           string    `json:"model"`
	ReasoningEffort string    `json:"reasoning_effort"`
	MaxTokens       int64     `json:"max_tokens"`
	Temperature     float64   `json:"temperature"`
	Thinking        *thinking `json:"thinking"`
	ResponseFormat  format    `json:"response_format"`
}

type message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

func (r *ChatRequestBody) AppendMessage(role, content string) {
	r.Messages = append(r.Messages, message{role, content})
}

type thinking struct {
	Type string `json:"type"`
}

func (r *ChatRequestBody) EnableThinking() {
	r.Thinking = &thinking{Type: constants.ThinkingEnabled}
}

func (r *ChatRequestBody) DisableThinking() {
	r.Thinking = nil
}

type format struct {
	Type string `json:"type"`
}

func (r *ChatRequestBody) SetResponseFormatText() {
	r.ResponseFormat = format{Type: constants.ResponseFormatText}
}

func (r *ChatRequestBody) SetResponseFormatJson() {
	r.ResponseFormat = format{Type: constants.ResponseFormatJSON}
}
