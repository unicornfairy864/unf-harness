package tool

import (
	"github.com/unicornfairy864/unf-harness/tool/internal/tools"
	"github.com/unicornfairy864/unf-harness/utils"
)

func Execute(functionTool tools.FunctionTool) string {
	res, err := functionTool.Execute()
	if err != nil {
		return "Tool \"" + functionTool.Name() + "\" failed: " + utils.Truncate(err.Error(), 256)
	}
	return res
}
