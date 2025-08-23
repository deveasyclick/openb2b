package types

type WebhookEvent struct {
	ID        string                 `json:"id"`
	Type      string                 `json:"type"`
	Data      map[string]interface{} `json:"data"`
	CreatedAt string                 `json:"created_at"`
	WebhookID string                 `json:"webhook_id"`
}
