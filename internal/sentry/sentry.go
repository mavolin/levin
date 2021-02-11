package sentry

import (
	"github.com/getsentry/sentry-go"
	"github.com/mavolin/adam/pkg/bot"
	"github.com/mavolin/adam/pkg/plugin"
	"github.com/mavolin/disstate/v3/pkg/state"

	"github.com/mavolin/levin/internal/config"
	"github.com/mavolin/levin/internal/meta"
)

var hubKey = "sentry.hub"

// Init initializes sentry using config.C.
func Init() error {
	return sentry.Init(sentry.ClientOptions{
		Dsn:              config.C.Sentry.DSN,
		Debug:            false,
		AttachStacktrace: false,
		SampleRate:       config.C.Sentry.SampleRate,
		ServerName:       config.C.ServerName,
		Release:          meta.Version,
		Environment:      config.C.Sentry.Environment,
	})
}

// Get extracts the *sentry.Hub from the passed context or event.
//
// If no hub is accessible using 'sentry.hub' as key, a clone of the current
// hub will be returned.
func Get(ctx interface{ Get(string) interface{} }) *sentry.Hub {
	if hub := ctx.Get(hubKey); hub != nil {
		if hub, ok := hub.(*sentry.Hub); ok && hub != nil {
			return hub
		}
	}

	return sentry.CurrentHub().Clone()
}

type Middlewares struct {
	MessageCreateMiddleware func(*state.State, *state.MessageCreateEvent)
	MessageUpdateMiddleware func(*state.State, *state.MessageUpdateEvent)

	BotMiddleware bot.MiddlewareFunc
}

// NewMiddlewares creates new collection of middlewares.
// The two MessageXMiddlewares store a clone of the passed *sentry.Hub in the
// events base.
//
// The BotMiddleware extracts that hub, and attaches meta information to its
// scope.
func NewMiddlewares(h *sentry.Hub) Middlewares {
	return Middlewares{
		MessageCreateMiddleware: func(_ *state.State, e *state.MessageCreateEvent) {
			h := h.Clone()
			h.Scope().SetExtra("message", e.Message)

			e.Set(hubKey, h)
		},
		MessageUpdateMiddleware: func(_ *state.State, e *state.MessageUpdateEvent) {
			h := h.Clone()
			h.Scope().SetExtra("message", e.Message)

			e.Set(hubKey, h)
		},
		BotMiddleware: func(next bot.CommandFunc) bot.CommandFunc {
			return func(s *state.State, ctx *plugin.Context) error {
				h := Get(ctx)
				h.Scope().SetTransaction(ctx.InvokedCommand.ProviderName + "/" + string(ctx.InvokedCommand.Identifier))
				h.Scope().SetTags(map[string]string{
					"err_source":      "bot",
					"plugin_provider": ctx.InvokedCommand.ProviderName,
					"command_id":      string(ctx.InvokedCommand.Identifier),
					"lang":            ctx.Lang,
				})

				return next(s, ctx)
			}
		},
	}
}
