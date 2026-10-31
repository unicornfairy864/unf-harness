package tool

import (
	"github.com/unicornfairy864/unf-harness/model"
)

var Tools map[string]model.Function

func Init() {
	Tools = make(map[string]model.Function)
	Tools[ToolReadFile.Name] = ToolReadFile
}
