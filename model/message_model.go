package model

import "encoding/json"

func (r *RequestBody) AppendRawMessage(msg json.RawMessage) {
	r.Messages = append(r.Messages, msg)
}

type CommonMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

func (r *RequestBody) AppendText(role, content string) {
	r.Messages = append(r.Messages, CommonMessage{Role: role, Content: content})
}

type ImageMessage struct {
	Role    string        `json:"role"`
	Content []interface{} `json:"content"`
}

func (r *RequestBody) AppendImage(data *string) {
	if data == nil {
		return
	}
	msg := &ImageMessage{
		Role: "user",
		Content: []interface{}{
			&struct {
				Type     string `json:"type"`
				ImageUrl struct {
					Url string `json:"url"`
				} `json:"image_url"`
			}{
				Type: "image_url",
				ImageUrl: struct {
					Url string `json:"url"`
				}{Url: *data},
			}},
	}
	r.Messages = append(r.Messages, msg)
}

type ToolMessage struct {
	Role      string     `json:"role"`
	Content   string     `json:"content"`
	ToolCalls []ToolCall `json:"tool_calls,omitempty"`
}
