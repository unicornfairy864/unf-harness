package config

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/viper"
)

// configFile 是唯一会读的配置文件，相对进程工作目录。
const configFile = "config.yaml"

// required 是 config.yaml 必须提供的 key；少任何一个都算"配置不完整"。
var required = []string{
	"llm.base_url", "llm.model", "llm.api_key", "llm.timeout", "llm.temperature", "llm.max_tokens", "llm.stream",
	"agent.max_turns", "agent.system_prompt", "agent.task_timeout", "agent.repeat_limit",
	"context.max_prompt_tokens", "context.keep_recent_turns", "context.tool_result_max_bytes",
	"tools.enabled", "tools.run_shell.timeout", "tools.run_shell.max_output_bytes", "tools.write_file.max_bytes",
	"trace.enabled", "trace.dir", "trace.record_payloads",
	"log.level",
}

// Init 只做两件事：把 config.yaml 读一次进 cfg，再用命令行 flag 覆盖。
// 不读环境变量、不写内置默认、不 watch、不建目录。
func Init(cfg *Config) {
	v := viper.New()
	v.SetConfigFile(configFile) // 固定路径：不 AddConfigPath、不写默认文件
	if err := v.ReadInConfig(); err != nil {
		panic(fmt.Errorf("config: 读取 %s 失败: %w", configFile, err))
	}
	if missing := missingKeys(v); len(missing) > 0 {
		panic(fmt.Errorf("config: %s 不完整，缺少: %s", configFile, strings.Join(missing, ", ")))
	}
	if err := v.Unmarshal(cfg); err != nil {
		panic(fmt.Errorf("config: 解析 %s 失败: %w", configFile, err))
	}

	args := new(Args)
	fs := flag.NewFlagSet("uh", flag.ExitOnError)
	args.bind(fs, cfg)
	_ = fs.Parse(os.Args[1:]) // -h 打 usage 退 0；参数写错打错误退 2
	args.apply(cfg)
}

// missingKeys 返回 config.yaml 里缺失的 key。
func missingKeys(v *viper.Viper) []string {
	var missing []string
	for _, k := range required {
		if !v.IsSet(k) {
			missing = append(missing, k)
		}
	}
	return missing
}
