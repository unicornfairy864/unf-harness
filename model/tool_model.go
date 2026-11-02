package model

type FunctionTool struct {
	Type     string   `json:"type"`
	Function Function `json:"function"`
}

func ToolToFuncTool(tool Function) FunctionTool {
	return FunctionTool{
		Type:     "function",
		Function: tool,
	}
}

type Function struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Strict      bool   `json:"strict"`

	Parameters Schema `json:"parameters"`
}

type Schema struct {
	Type     string   `json:"type"` // "object"
	Required []string `json:"required"`

	Properties map[string]Property `json:"properties"`
}

type Property struct {
	Type        string   `json:"type"`
	Description string   `json:"description"`
	Enum        []string `json:"enum,omitempty"`
}
