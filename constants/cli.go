package constants

const (
	CliExitCommand      = "exit"
	CliExitPrompt       = "Goodbye!"
	CliInputPrompt      = "User > "
	CliInputAgainPrompt = "Please input again. | or Type \"" + CliExitPrompt + "\" to quit."
)

var PromptStop = map[string]string{
	StopReasonStop:          "System > Success.",
	StopReasonLength:        "System > Reached the maximum token limit.",
	StopReasonToolCalls:     "System > Executing tool calls...",
	StopReasonContentFilter: "System > Content was filtered by remote server.",
}
