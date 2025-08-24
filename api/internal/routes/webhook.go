package routes

import (
	"github.com/deveasyclick/openb2b/internal/modules/webhook"
	"github.com/deveasyclick/openb2b/internal/shared/deps"
	"github.com/deveasyclick/openb2b/pkg/interfaces"
	"github.com/go-chi/chi"
)

func registerWebhookRoutes(r chi.Router, webhookHandler interfaces.WebhookHandler, appCtx *deps.AppContext) {

	r.Route("/webhooks", func(r chi.Router) {

		r.With(webhook.Verify(appCtx)).
			Post("/createUser", webhookHandler.HandleClerkEvents)
	})
}
