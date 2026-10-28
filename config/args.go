package config

import (
	"flag"
	"time"
)

// Args 是唯一允许从命令行覆盖的配置（llm 段）。
type Args struct {
	Model       string
	ApiKey      string
	TimeOut     time.Duration
	Temperature float64
	MaxTokens   int64
	Stream      bool
}

// bind 注册 flag：默认值取刚读出来的 yaml 值 → 不传=用 yaml，传了=覆盖。
func (a *Args) bind(fs *flag.FlagSet, cfg *Config) {
	fs.StringVar(&a.Model, "model", cfg.LLM.Model, "覆盖 llm.model")
	fs.StringVar(&a.ApiKey, "api-key", cfg.LLM.APIKey, "覆盖 llm.api_key")
	fs.DurationVar(&a.TimeOut, "timeout", cfg.LLM.Timeout, "覆盖 llm.timeout")
	fs.Float64Var(&a.Temperature, "temperature", cfg.LLM.Temperature, "覆盖 llm.temperature")
	fs.Int64Var(&a.MaxTokens, "max-tokens", cfg.LLM.MaxTokens, "覆盖 llm.max_tokens")
	fs.BoolVar(&a.Stream, "stream", cfg.LLM.Stream, "覆盖 llm.stream")
}

// apply 把 flag 的值直接写回 Config（无条件覆盖）。
func (a *Args) apply(cfg *Config) {
	cfg.LLM.Model = a.Model
	cfg.LLM.APIKey = a.ApiKey
	cfg.LLM.Timeout = a.TimeOut
	cfg.LLM.Temperature = a.Temperature
	cfg.LLM.MaxTokens = a.MaxTokens
	cfg.LLM.Stream = a.Stream
}
