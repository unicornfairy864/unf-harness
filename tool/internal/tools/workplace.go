package tools

import (
	"path/filepath"

	"github.com/unicornfairy864/unf-harness/config"
)

// workplaceRoot 返回配置里的工作目录，没配就退化成当前目录。
func workplaceRoot() string {
	if root := config.UH_Config.Tools.Workplace; root != "" {
		return root
	}
	return "."
}

// resolvePath 把模型给的路径落成实际路径：
// 「大写字母:」开头当绝对路径，其余一律挂到 workplace 下。
func resolvePath(p string) string {
	if isDriveAbsPath(p) {
		return p
	}
	return filepath.Join(workplaceRoot(), p)
}
