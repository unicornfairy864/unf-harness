package tools

import (
	"github.com/unicornfairy864/unf-harness/model"
)

type FunctionTool interface {
	Name() string
	ReceiveParam(message string) error
	Execute() (string, error)
	GetTool() model.Function
}
