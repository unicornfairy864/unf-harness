package constants

const (
	CliExitCommand      = "exit"
	CliExitPrompt       = "\nGoodbye!"
	CliUserPrompt       = "\nUser > "
	CliContPrompt       = "     > "
	CliThinkingPrompt   = "\n● Thinking\n"
	CliAssistantPrompt  = "\n● Assistant\n"
	CliInputAgainPrompt = "Please input again. | or Type \"" + CliExitCommand + "\" to quit."
)

var PromptStop = map[string]string{
	StopReasonStop:          "[System]Success.",
	StopReasonLength:        "[System]Reached the maximum token limit.",
	StopReasonToolCalls:     "[System]Executing tool calls...",
	StopReasonContentFilter: "[System]Content was filtered by remote server.",
}
