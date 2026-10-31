package model

import (
	"github.com/unicornfairy864/unf-harness/constants"
)

type RequestBody struct {
	Model           string  `json:"model"`
	ReasoningEffort string  `json:"reasoning_effort"`
	MaxTokens       int64   `json:"max_tokens"`
	Temperature     float64 `json:"temperature"`

	Thinking       *thinking `json:"thinking"`
	ResponseFormat format    `json:"response_format"`

	Messages []interface{} `json:"messages"`

	Tools []FunctionTool `json:"tools,omitempty"`
}

type thinking struct {
	Type string `json:"type"`
}

func (r *RequestBody) EnableThinking() {
	r.Thinking = &thinking{Type: constants.ThinkingEnabled}
}

func (r *RequestBody) DisableThinking() {
	r.Thinking = &thinking{Type: constants.ThinkingDisabled}
}

type format struct {
	Type string `json:"type"`
}

func (r *RequestBody) SetResponseFormatText() {
	r.ResponseFormat = format{Type: constants.ResponseFormatText}
}

func (r *RequestBody) SetResponseFormatJson() {
	r.ResponseFormat = format{Type: constants.ResponseFormatJSON}
}

func (r *RequestBody) RegisterTool(t Function) {
	r.Tools = append(r.Tools, ToolToFuncTool(t))
}
