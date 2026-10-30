package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/unicornfairy864/unf-harness/client/internal/models"
	"github.com/unicornfairy864/unf-harness/config"
)

func (c *LLMClient) Request(postBody models.ChatRequestBody) (*models.Data, error) {
	// 设置 body 请求
	dataStr, _ := json.Marshal(postBody)

	req, _ := http.NewRequest("POST", config.UH_Config.LLM.BaseURL,
		bytes.NewReader(dataStr))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+config.UH_Config.LLM.APIKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 4<<10))
		return nil, fmt.Errorf("unexpected status %d: %s", resp.StatusCode, body)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	result, err := parseResp(body)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func parseResp(body []byte) (*models.Data, error) {
	/*  现在只能解析 deepseek
	 *  以后可能兼容更多服务提供商
	 *  为什么不用 openai-go: 因为 esp32 也要用 http 请求
	 */
	data := models.Data{}
	err := json.Unmarshal(body, &data)
	if err != nil {
		return nil, err
	}
	return &data, nil
}
