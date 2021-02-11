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
		TracesSampleRate: config.C.Sentry.TracesSampleRate,
		ServerName:       config.C.ServerName,
		Release:          meta.Version,
		Environment:      config.C.Sentry.Environment,
	})
}

// GetHub extracts the *sentry.Hub from the passed *plugin.Context.
//
// If no hub is accessible using 'sentry.hub' as key, a clone of the current
// hub will be returned.
func GetHub(ctx *plugin.Context) *sentry.Hub {
	if hub := ctx.Get(hubKey); hub != nil {
		if hub, ok := hub.(*sentry.Hub); ok && hub != nil {
			return hub
		}
	}

	return sentry.CurrentHub().Clone()
}

// NewMiddleware creates a new bot.MiddlewareFunc from the passed *sentry.Hub,
// that stores a *sentry.Hub in the command's context under the key
// 'sentry.hub'.
func NewMiddleware(h *sentry.Hub) bot.MiddlewareFunc {
	return func(next bot.CommandFunc) bot.CommandFunc {
		return func(s *state.State, ctx *plugin.Context) error {
			h := h.Clone()

			h.Scope().SetTransaction(ctx.InvokedCommand.ProviderName + "/" + string(ctx.InvokedCommand.Identifier))
			h.Scope().SetTags(map[string]string{
				"err_source":      "bot",
				"plugin_provider": ctx.InvokedCommand.ProviderName,
				"command_id":      string(ctx.InvokedCommand.Identifier),
				"lang":            ctx.Lang,
			})
			h.Scope().SetExtra("message", ctx.Message)

			ctx.Set(hubKey, h)
			return next(s, ctx)
		}
	}
}
