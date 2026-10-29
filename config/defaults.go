package config

import "time"

// 第一层：程序内置默认值，逐个对应 config.yaml 的 key。
// 这是最底层的兜底：config.yaml 没写、命令行也没传的 key，最终就是这些值。
// 改这里 = 改 config.yaml 缺失时的行为；config.yaml 和 flag 都能覆盖它。
const (
	// llm
	defaultLLMBaseURL     = "http://127.0.0.1:8080/v1" // 换远程 API 只改这一行
	defaultLLMModel       = "qwen3-coder"
	defaultLLMAPIKey      = "" // 空：本地 llama-server 不校验，远程用 -k 传
	defaultLLMTimeout     = 5 * time.Minute
	defaultLLMTemperature = 0.2
	defaultLLMMaxTokens   = 4096
	defaultLLMStream      = false // MVP 关掉

	// agent
	defaultAgentMaxTurns     = 25
	defaultAgentSystemPrompt = "You are a helpful assistant."
	defaultAgentUserPrompt   = "" // 空 = validate 阶段 panic，必须由 config.yaml 或 -c 给出
	defaultAgentTaskTimeout  = 60 * time.Second
	defaultAgentRepeatLimit  = 3

	// context
	defaultContextMaxPromptTokens    = 24000 // 必须 < llama-server 的 -c 32768，给生成留空间
	defaultContextKeepRecentTurns    = 6
	defaultContextToolResultMaxBytes = 8192

	// tools
	defaultRunShellTimeout        = 60 * time.Second
	defaultRunShellMaxOutputBytes = 16384
	defaultWriteFileMaxBytes      = 1048576

	// trace
	defaultTraceEnabled        = true
	defaultTraceDir            = "./traces"
	defaultTraceRecordPayloads = true

	// log
	defaultLogLevel = "info"
)

// defaultToolsEnabled 只能是变量：切片不是常量表达式。
var defaultToolsEnabled = []string{"run_shell", "read_file", "write_file"}

// defaultConfig 返回第一层（const）填满的 Config，作为整个配置合并的起点。
func defaultConfig() Config {
	return Config{
		LLM: LLM{
			BaseURL:     defaultLLMBaseURL,
			Model:       defaultLLMModel,
			APIKey:      defaultLLMAPIKey,
			Timeout:     defaultLLMTimeout,
			Temperature: defaultLLMTemperature,
			MaxTokens:   defaultLLMMaxTokens,
			Stream:      defaultLLMStream,
		},
		Agent: Agent{
			MaxTurns:     defaultAgentMaxTurns,
			SystemPrompt: defaultAgentSystemPrompt,
			UserPrompt:   defaultAgentUserPrompt,
			TaskTimeout:  defaultAgentTaskTimeout,
			RepeatLimit:  defaultAgentRepeatLimit,
		},
		Context: Context{
			MaxPromptTokens:    defaultContextMaxPromptTokens,
			KeepRecentTurns:    defaultContextKeepRecentTurns,
			ToolResultMaxBytes: defaultContextToolResultMaxBytes,
		},
		Tools: Tools{
			Enabled: append([]string(nil), defaultToolsEnabled...), // 拷一份，别让上层改到常量层
			RunShell: RunShell{
				Timeout:        defaultRunShellTimeout,
				MaxOutputBytes: defaultRunShellMaxOutputBytes,
			},
			WriteFile: WriteFile{
				MaxBytes: defaultWriteFileMaxBytes,
			},
		},
		Trace: Trace{
			Enabled:        defaultTraceEnabled,
			Dir:            defaultTraceDir,
			RecordPayloads: defaultTraceRecordPayloads,
		},
		Log: Log{
			Level: defaultLogLevel,
		},
	}
}
