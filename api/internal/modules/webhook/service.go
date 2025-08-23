package webhook

import (
	"context"
	"fmt"
	"net/http"

	"github.com/deveasyclick/openb2b/internal/model"
	"github.com/deveasyclick/openb2b/internal/shared/apperrors"
	"github.com/deveasyclick/openb2b/internal/shared/deps"
	"github.com/deveasyclick/openb2b/internal/shared/types"
	"github.com/deveasyclick/openb2b/pkg/interfaces"
	"github.com/mitchellh/mapstructure"
)

type webhookService struct {
	userService   interfaces.UserService
	clerkService  interfaces.ClerkService
	appCtx        *deps.AppContext
	eventHandlers map[string]func(context.Context, map[string]interface{}) *apperrors.APIError
}

func NewWebhookService(service interfaces.UserService, clerkService interfaces.ClerkService, appCtx *deps.AppContext) interfaces.WebhookService {
	s := &webhookService{userService: service, clerkService: clerkService, appCtx: appCtx}

	// eventHandlers maps webhook event types to their corresponding handler functions.
	//
	// Each entry in this map associates a Clerk webhook event type (e.g. "user.created")
	// with a method on the `webhookService` that knows how to handle it. The handlers
	// follow a common function signature:
	//
	//   func(ctx context.Context, data map[string]interface{}) *apperrors.APIError
	//
	// This makes it easy to add new webhook event support by simply adding another
	// entry to the map, without modifying the routing logic in HandleEvent.
	//
	// Example:
	//
	//   var eventHandlers = map[string]func(context.Context, map[string]interface{}) *apperrors.APIError{
	//       "user.created": s.CreateUser,
	//       "user.updated": s.UpdateUser,
	//   }
	//
	// This design keeps webhook handling modular and extensible.
	s.eventHandlers = map[string]func(context.Context, map[string]interface{}) *apperrors.APIError{
		"user.created": s.createUser,
	}

	return s
}

// HandleEvent routes an incoming webhook event to the appropriate handler.
//
// It looks up the event type in the service's `eventHandlers` map and, if a
// matching handler is found, executes it with the provided context and event
// data. The handler is expected to return an *apperrors.APIError if any
// application-level error occurs.
//
// If no handler is registered for the event type, the method logs the event
// as ignored and returns nil (indicating no error).
//
// Parameters:
//   - ctx context.Context: request-scoped context passed through to the handler.
//   - event *types.WebhookEvent: the parsed webhook event containing the type
//     and payload data.
//
// Returns:
//   - *apperrors.APIError: a structured error if the handler fails, or nil if
//     the event was processed successfully or ignored.
//
// Example:
//
//   err := webhookService.HandleEvent(ctx, event)
//   if err != nil {
//       response.WriteJSONError(w, err, appCtx.Logger)
//       return
//   }
//
// This allows the webhook service to be easily extended by registering new
// event handlers in the `eventHandlers` map.

func (s *webhookService) HandleEvent(ctx context.Context, event *types.WebhookEvent) *apperrors.APIError {
	if handler, ok := s.eventHandlers[event.Type]; ok {
		return handler(ctx, event.Data)
	}

	s.appCtx.Logger.Info("Ignoring unknown webhook event type", "type", event.Type)
	return nil
}

func (s *webhookService) createUser(ctx context.Context, data map[string]interface{}) *apperrors.APIError {
	var userData types.ClerkUser
	if err := mapstructure.Decode(data, &userData); err != nil {
		return &apperrors.APIError{
			Code:    http.StatusBadRequest,
			Message: fmt.Sprintf("failed to decode ClerkUser: %s", err),
		}
	}

	email := ""
	if len(userData.EmailAddresses) > 0 {
		email = userData.EmailAddresses[0].EmailAddress
	} else {
		return &apperrors.APIError{
			Code:    http.StatusBadRequest,
			Message: "no email address found in ClerkUser",
		}
	}

	user := &model.User{
		ClerkID:   userData.ID,
		FirstName: userData.FirstName,
		LastName:  userData.LastName,
		Email:     email,
		Role:      string(model.RoleAdmin),
	}

	err := s.userService.Create(ctx, user)
	if err != nil {
		return &apperrors.APIError{
			Code:    http.StatusInternalServerError,
			Message: "error creating distributor in webhook",
		}
	}

	err = s.clerkService.SetExternalID(user)
	if err != nil {
		s.appCtx.Logger.Error("error updating clerk user", "error", err, "user", user.ID)
	} else {
		s.appCtx.Logger.Info("Updated clerk user externalId", "externalId", user.ID)
	}

	return nil
}
