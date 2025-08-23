package svix

import (
	"sync"

	svix "github.com/svix/svix-webhooks/go"
)

var (
	wh     *svix.Webhook
	whOnce sync.Once
	whErr  error
)

// GetWebhookVerifier returns a singleton instance of the Svix webhook verifier
func GetWebhookVerifier(secret string) (*svix.Webhook, error) {
	whOnce.Do(func() {
		wh, whErr = svix.NewWebhook(secret)
	})

	return wh, whErr
}
