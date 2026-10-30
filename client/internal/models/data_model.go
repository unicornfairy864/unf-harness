package models

type Data struct {
	Id                string `json:"id"`
	Model             string `json:"model"`
	Created           int    `json:"created"`
	Object            string `json:"object"`
	SystemFingerprint string `json:"system_fingerprint"`

	Choices []struct {
		FinishReason string  `json:"finish_reason"`
		Index        int     `json:"index"`
		Logprobs     *string `json:"logprobs"`
		Message      struct {
			Content          string `json:"content"`
			ReasoningContent string `json:"reasoning_content"`
			Role             string `json:"role"`
		} `json:"message"`
	} `json:"choices"`

	Usage struct {
		CompletionTokens        int `json:"completion_tokens"`
		CompletionTokensDetails struct {
			ReasoningTokens int `json:"reasoning_tokens"`
		} `json:"completion_tokens_details"`
		PromptCacheHitTokens  int `json:"prompt_cache_hit_tokens"`
		PromptCacheMissTokens int `json:"prompt_cache_miss_tokens"`
		PromptTokens          int `json:"prompt_tokens"`
		PromptTokensDetails   struct {
			CachedTokens int `json:"cached_tokens"`
		} `json:"prompt_tokens_details"`
		TotalTokens int `json:"total_tokens"`
	} `json:"usage"`
}
