package api

import "encoding/json"

// WebhookConfig 服务端返回的 Webhook 配置
type WebhookConfig struct {
	URL        string `json:"url"`
	IsVerified bool   `json:"isVerified"`
	RetryCount int    `json:"retryCount"`
	CreatedAt  string `json:"createdAt"`
}

// WebhookSetRequest 设置 Webhook 的请求体（CLI 通过 apikey 认证）
type WebhookSetRequest struct {
	Apikey     string `json:"apikey"`
	Role       string `json:"role,omitempty"`
	URL        string `json:"url"`
	Secret     string `json:"secret,omitempty"`
	RetryCount int    `json:"retryCount,omitempty"`
}

// GetWebhook 获取当前 Webhook 配置
func (c *Client) GetWebhook(apikey string) (*WebhookConfig, error) {
	resp, err := c.Get("/v1/cli/webhook?apikey=" + apikey)
	if err != nil {
		return nil, err
	}
	// data 可能为 null（未配置）
	if string(resp.Data) == "null" {
		return nil, nil
	}
	var cfg WebhookConfig
	if err := json.Unmarshal(resp.Data, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

// SetWebhook 设置或更新 Webhook 配置
func (c *Client) SetWebhook(req *WebhookSetRequest) error {
	_, err := c.Put("/v1/cli/webhook", req)
	return err
}

// RemoveWebhook 移除 Webhook 配置
func (c *Client) RemoveWebhook(apikey string) error {
	_, err := c.Delete("/v1/cli/webhook", map[string]string{"apikey": apikey})
	return err
}

// TestWebhook 发送测试回调验证连通性
func (c *Client) TestWebhook(apikey string) error {
	_, err := c.Post("/v1/cli/webhook/test", map[string]string{"apikey": apikey})
	return err
}
