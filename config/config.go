package config

import "time"

// Config 与 config.yaml 一一对应：section -> field。
type Config struct {
	LLM   LLM   `mapstructure:"llm"`
	Tools Tools `mapstructure:"tools"`
	Trace Trace `mapstructure:"trace"`
	Log   Log   `mapstructure:"log"`
}

type LLM struct {
	BaseURL            string        `mapstructure:"base_url"`
	Model              string        `mapstructure:"model"`
	APIKey             string        `mapstructure:"api_key"`
	Timeout            time.Duration `mapstructure:"timeout"`
	Temperature        float64       `mapstructure:"temperature"`
	MaxTokens          int64         `mapstructure:"max_tokens"`
	Stream             bool          `mapstructure:"stream"`
	Thinking           bool          `mapstructure:"thinking"`
	ReasoningEffort    string        `mapstructure:"reasoning_effort"`
	MaxTurns           int           `mapstructure:"max_turns"`
	SystemPrompt       string        `mapstructure:"system_prompt"`
	TaskTimeout        time.Duration `mapstructure:"task_timeout"`
	RepeatLimit        int           `mapstructure:"repeat_limit"`
	MaxPromptTokens    int           `mapstructure:"max_prompt_tokens"`
	KeepRecentTurns    int           `mapstructure:"keep_recent_turns"`
	ToolResultMaxBytes int           `mapstructure:"tool_result_max_bytes"`
}

type Tools struct {
	Enabled   []string  `mapstructure:"enabled"`
	RunShell  RunShell  `mapstructure:"run_shell"`
	WriteFile WriteFile `mapstructure:"write_file"`
}

type RunShell struct {
	Timeout        time.Duration `mapstructure:"timeout"`
	MaxOutputBytes int           `mapstructure:"max_output_bytes"`
}

type WriteFile struct {
	MaxBytes int `mapstructure:"max_bytes"`
}

type Trace struct {
	Enabled        bool   `mapstructure:"enabled"`
	Dir            string `mapstructure:"dir"`
	RecordPayloads bool   `mapstructure:"record_payloads"`
}

type Log struct {
	Level string `mapstructure:"level"`
}
