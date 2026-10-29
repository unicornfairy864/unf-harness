package config

import (
	"errors"
	"flag"
	"fmt"
	"io/fs"
	"os"

	"github.com/spf13/viper"
)

// configFile 是第二层唯一会读的配置文件，相对进程工作目录。
const configFile = "config.yaml"

// Init 按三层顺序合并配置：程序 const → config.yaml → 命令行 flag。
// 每一层只覆盖它显式给出的 key，没给的一律沿用上一层。
// 合并完做一次校验，配置不合法直接 panic（别带着坏配置往下跑）。
func Init(cfg *Config) {
	*cfg = defaultConfig() // 第一层：程序内置默认值
	loadYAML(cfg)          // 第二层：config.yaml
	applyArgs(cfg)         // 第三层：命令行 flag
	if err := cfg.validate(); err != nil {
		panic(fmt.Errorf("config: %w", err))
	}
}

// loadYAML 把 config.yaml 里出现的 key 覆盖到 cfg 上。
// 文件不存在不算错误：整份配置退化成只有 const 层，只在 stderr 提示一句；
// 文件存在但读不动/写坏了才 panic（那种错不提示就直接跑更危险）。
func loadYAML(cfg *Config) {
	v := viper.New()
	v.SetConfigFile(configFile) // 固定路径：不 AddConfigPath、不搜目录、不 watch

	if err := v.ReadInConfig(); err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			fmt.Fprintf(os.Stderr, "config: 未找到 %s，全部使用程序内置默认值\n", configFile)
			return
		}
		panic(fmt.Errorf("config: 读取 %s 失败: %w", configFile, err))
	}

	// Unmarshal 只写 yaml 里出现过的 key，缺的字段保持 const 层的值。
	if err := v.Unmarshal(cfg); err != nil {
		panic(fmt.Errorf("config: 解析 %s 失败: %w", configFile, err))
	}
}

// applyArgs 解析命令行 flag（第三层），默认值取前两层合并后的值。
func applyArgs(cfg *Config) {
	args := new(Args)
	set := flag.NewFlagSet("uh", flag.ExitOnError)
	args.bind(set, cfg)
	_ = set.Parse(os.Args[1:]) // -h/-help 打 usage 退 0；参数写错打错误退 2
	args.apply(cfg)
}
