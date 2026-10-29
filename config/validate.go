package config

import (
	"fmt"
	"strings"
)

// validLogLevels 是 log.level 的白名单。
var validLogLevels = map[string]struct{}{
	"debug": {},
	"info":  {},
	"warn":  {},
	"error": {},
}

// validate 在三层合并完之后做一次体检：任何一处不合法就返回错误，调用方直接 panic。
// 只查"错了就一定跑不起来"的项，不替上层做业务判断。
func (c *Config) validate() error {
	var bad []string
	bad = append(bad, c.LLM.problems()...)
	bad = append(bad, c.Agent.problems()...)
	bad = append(bad, c.Context.problems()...)
	bad = append(bad, c.Tools.problems()...)
	bad = append(bad, c.Trace.problems()...)
	bad = append(bad, c.Log.problems()...)
	if len(bad) > 0 {
		return fmt.Errorf("%s 配置不合法: %s", configFile, strings.Join(bad, "; "))
	}
	return nil
}

func (l LLM) problems() []string {
	var bad []string
	if l.BaseURL == "" {
		bad = append(bad, "llm.base_url 不能为空")
	}
	if l.Model == "" {
		bad = append(bad, "llm.model 不能为空")
	}
	if l.Timeout <= 0 {
		bad = append(bad, fmt.Sprintf("llm.timeout 必须 > 0（当前 %s）", l.Timeout))
	}
	if l.Temperature < 0 || l.Temperature > 2 {
		bad = append(bad, fmt.Sprintf("llm.temperature 必须落在 [0,2]（当前 %g）", l.Temperature))
	}
	if l.MaxTokens <= 0 {
		bad = append(bad, fmt.Sprintf("llm.max_tokens 必须 > 0（当前 %d）", l.MaxTokens))
	}
	// 只有开着思考才会把 reasoning_effort 发出去，关掉时留空不算错。
	if l.Thinking && l.ReasoningEffort == "" {
		bad = append(bad, "llm.thinking 开启时 llm.reasoning_effort 不能为空（如 low / medium / high）")
	}
	return bad
}

func (a Agent) problems() []string {
	var bad []string
	if a.MaxTurns <= 0 {
		bad = append(bad, fmt.Sprintf("agent.max_turns 必须 > 0（当前 %d）", a.MaxTurns))
	}
	if a.SystemPrompt == "" {
		bad = append(bad, "agent.system_prompt 不能为空")
	}
	if a.UserPrompt == "" {
		bad = append(bad, "agent.user_prompt 不能为空（config.yaml 里留空就必须用 -c 传）")
	}
	if a.TaskTimeout <= 0 {
		bad = append(bad, fmt.Sprintf("agent.task_timeout 必须 > 0（当前 %s）", a.TaskTimeout))
	}
	if a.RepeatLimit <= 0 {
		bad = append(bad, fmt.Sprintf("agent.repeat_limit 必须 > 0（当前 %d）", a.RepeatLimit))
	}
	return bad
}

func (c Context) problems() []string {
	var bad []string
	if c.MaxPromptTokens <= 0 {
		bad = append(bad, fmt.Sprintf("context.max_prompt_tokens 必须 > 0（当前 %d）", c.MaxPromptTokens))
	}
	if c.KeepRecentTurns < 0 {
		bad = append(bad, fmt.Sprintf("context.keep_recent_turns 不能为负（当前 %d）", c.KeepRecentTurns))
	}
	if c.ToolResultMaxBytes <= 0 {
		bad = append(bad, fmt.Sprintf("context.tool_result_max_bytes 必须 > 0（当前 %d）", c.ToolResultMaxBytes))
	}
	return bad
}

func (t Tools) problems() []string {
	var bad []string
	if len(t.Enabled) == 0 {
		bad = append(bad, "tools.enabled 不能为空（没有工具就不是 harness 了）")
	}
	if t.RunShell.Timeout <= 0 {
		bad = append(bad, fmt.Sprintf("tools.run_shell.timeout 必须 > 0（当前 %s）", t.RunShell.Timeout))
	}
	if t.RunShell.MaxOutputBytes <= 0 {
		bad = append(bad, fmt.Sprintf("tools.run_shell.max_output_bytes 必须 > 0（当前 %d）", t.RunShell.MaxOutputBytes))
	}
	if t.WriteFile.MaxBytes <= 0 {
		bad = append(bad, fmt.Sprintf("tools.write_file.max_bytes 必须 > 0（当前 %d）", t.WriteFile.MaxBytes))
	}
	return bad
}

func (t Trace) problems() []string {
	if t.Enabled && t.Dir == "" {
		return []string{"trace.enabled 为 true 时 trace.dir 不能为空"}
	}
	return nil
}

func (l Log) problems() []string {
	if _, ok := validLogLevels[l.Level]; !ok {
		return []string{fmt.Sprintf("log.level 必须是 debug/info/warn/error 之一（当前 %q）", l.Level)}
	}
	return nil
}
