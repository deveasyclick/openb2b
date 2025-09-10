package webhook

import (
	"context"
	"errors"
	"strconv"

	"github.com/deveasyclick/openb2b/internal/model"
	"github.com/deveasyclick/openb2b/internal/shared/apperrors"
	"github.com/deveasyclick/openb2b/internal/shared/deps"
	"github.com/deveasyclick/openb2b/internal/shared/types"
	"github.com/deveasyclick/openb2b/pkg/clerk"
	"github.com/deveasyclick/openb2b/pkg/interfaces"
	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

type service struct {
	userService   interfaces.UserService
	clerkService  clerk.Service
	appCtx        *deps.AppContext
	eventHandlers map[string]func(context.Context, map[string]interface{}) error
}

func NewService(us interfaces.UserService, cs clerk.Service, appCtx *deps.AppContext) interfaces.WebhookService {
	s := &service{userService: us, clerkService: cs, appCtx: appCtx}

	// eventHandlers maps webhook event types to their corresponding handler functions.
	//
	// Each entry in this map associates a Clerk webhook event type (e.g. "user.created")
	// with a method on the `service` that knows how to handle it. The handlers
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
	s.eventHandlers = map[string]func(context.Context, map[string]interface{}) error{
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
//   err := service.HandleEvent(ctx, event)
//   if err != nil {
//       response.WriteJSONError(w, err, appCtx.Logger)
//       return
//   }
//
// This allows the webhook service to be easily extended by registering new
// event handlers in the `eventHandlers` map.

func (s *service) HandleEvent(ctx context.Context, event *types.WebhookEvent) error {
	if handler, ok := s.eventHandlers[event.Type]; ok {
		return handler(ctx, event.Data)
	}

	s.appCtx.Logger.Info("Ignoring unknown webhook event type", "type", event.Type)
	return nil
}

func (s *service) createUser(ctx context.Context, data map[string]interface{}) error {
	var userData types.ClerkUser
	if err := mapstructure.Decode(data, &userData); err != nil {
		return errors.New(apperrors.ErrDecodeRequestBody)
	}

	email := ""
	if len(userData.EmailAddresses) > 0 {
		email = userData.EmailAddresses[0].EmailAddress
	} else {
		return errors.New(apperrors.ErrEmailNotFoundInClerkWebhook)
	}
	user, err := s.userService.FindByEmail(ctx, email)
	s.appCtx.Logger.Info("user", err)
	if err != nil && err != gorm.ErrRecordNotFound {
		return err
	}

	if user != nil {
		return errors.New(apperrors.ErrUserAlreadyExists)
	}

	user = &model.User{
		ClerkID:   userData.ID,
		FirstName: userData.FirstName,
		LastName:  userData.LastName,
		Email:     email,
		Role:      model.RoleAdmin,
	}

	err = s.userService.Create(ctx, user)
	if err != nil {
		// TODO: Move this later to a background worker
		go func(ctx context.Context, clerkId string) {
			cleanupError := s.clerkService.DeleteUser(ctx, user.ClerkID)
			if cleanupError != nil {
				s.appCtx.Logger.Error("error deleting user from clerk", "error", cleanupError, "user", user.ID)
			}
		}(ctx, user.ClerkID)
		return err
	}
	s.appCtx.Logger.Info("User created", "user ID", user.ID)

	userId := strconv.FormatUint(uint64(user.ID), 10)
	err = s.clerkService.SetRoleAndExternalID(ctx, user.ClerkID, userId, model.RoleAdmin)
	if err != nil {
		return err
	}
	s.appCtx.Logger.Info("Updated clerk user externalId", "externalId", user.ID)

	return nil
}
