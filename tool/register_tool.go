package tool

import (
	"github.com/unicornfairy864/unf-harness/tool/internal/tools"
)

var Tools map[string]tools.FunctionTool

func Register() {
	Tools = make(map[string]tools.FunctionTool)
	Tools[tools.ToolExample.Name()] = &tools.ToolExample
	Tools[tools.ToolReadFile.Name()] = &tools.ToolReadFile
	Tools[tools.ToolWriteFile.Name()] = &tools.ToolWriteFile
	Tools[tools.ToolRunShell.Name()] = &tools.ToolRunShell
}
