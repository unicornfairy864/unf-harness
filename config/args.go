package config

import (
	"flag"
	"time"
)

// Args 是唯一允许从命令行覆盖的配置（llm 段 + agent 的两个 prompt）。
type Args struct {
	Model           string        // -m
	MaxTokens       int64         // -M
	ApiKey          string        // -k
	SystemPrompt    string        // -S
	UserPrompt      string        // -c
	TimeOut         time.Duration // -t
	Temperature     float64       // -T
	Stream          bool          // -s
	Thinking        bool          // -thinking
	ReasoningEffort string        // -r
}

// bind 注册 flag：默认值取前两层合并后的值 → 不传=沿用，传了=覆盖。
func (a *Args) bind(fs *flag.FlagSet, cfg *Config) {
	fs.StringVar(&a.Model, "m", cfg.LLM.Model, "选择使用的模型")
	fs.Int64Var(&a.MaxTokens, "M", cfg.LLM.MaxTokens, "最大输出token数")
	fs.StringVar(&a.ApiKey, "k", cfg.LLM.APIKey, "api_key")
	fs.StringVar(&a.SystemPrompt, "S", cfg.Agent.SystemPrompt, "system_prompt")
	fs.StringVar(&a.UserPrompt, "c", cfg.Agent.UserPrompt, "user_prompt")
	fs.DurationVar(&a.TimeOut, "t", cfg.LLM.Timeout, "最大输出时间")
	fs.Float64Var(&a.Temperature, "T", cfg.LLM.Temperature, "temperature")
	fs.BoolVar(&a.Stream, "s", cfg.LLM.Stream, "是否启用流式输出")
	fs.BoolVar(&a.Thinking, "thinking", cfg.LLM.Thinking, "是否开启思考（关闭请用 -thinking=false）")
	fs.StringVar(&a.ReasoningEffort, "r", cfg.LLM.ReasoningEffort, "reasoning_effort")
}

// apply 把 flag 的值写回 Config（第三层，无条件覆盖）。
func (a *Args) apply(cfg *Config) {
	cfg.LLM.Model = a.Model
	cfg.LLM.MaxTokens = a.MaxTokens
	cfg.LLM.APIKey = a.ApiKey
	cfg.LLM.Timeout = a.TimeOut
	cfg.LLM.Temperature = a.Temperature
	cfg.LLM.Stream = a.Stream
	cfg.LLM.Thinking = a.Thinking
	cfg.LLM.ReasoningEffort = a.ReasoningEffort
	cfg.Agent.SystemPrompt = a.SystemPrompt
	cfg.Agent.UserPrompt = a.UserPrompt
}
