package tool

import (
	"github.com/unicornfairy864/unf-harness/tool/internal/tools"
)

func Execute(functionTool tools.FunctionTool) string {
	res, err := functionTool.Execute()
	if err != nil {
		return "Tool \"" + functionTool.Name() + "\" failed: " + err.Error()
	}
	return res
}
