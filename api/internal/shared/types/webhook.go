package types

type contextKey string

const WebhookEventKey contextKey = "webhookEvent"

// WebhookEvent represents a Clerk webhook payload
// swagger:model WebhookEvent
type WebhookEvent struct {
	// ID of the webhook event
	ID string `json:"id" example:"evt_123456"`

	// Type of the webhook event (e.g., "user.created")
	Type string `json:"type" example:"user.created"`

	// Data payload of the webhook
	Data map[string]interface{} `json:"data"`

	// Event creation timestamp
	CreatedAt string `json:"created_at" example:"2025-08-26T18:00:00Z"`

	// ID of the webhook that triggered this event
	WebhookID string `json:"webhook_id" example:"wh_987654"`
}
