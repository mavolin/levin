package sentry

import (
	"github.com/getsentry/sentry-go"
	"github.com/mavolin/adam/pkg/bot"
	"github.com/mavolin/adam/pkg/plugin"
	"github.com/mavolin/disstate/v3/pkg/state"

	"github.com/mavolin/levin/internal/config"
	"github.com/mavolin/levin/internal/meta"
)

const hubKey = "sentry.hub"

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

// Get extracts the *sentry.Hub from the *plugin.Context.
// This will only work, if a *sentry.Hub was previously stored in the Context
// under 'sentry.hub', e.g. using NewMiddleware.
func Get(ctx *plugin.Context) *sentry.Hub {
	return ctx.Get(hubKey).(*sentry.Hub)
}

// NewMiddleware creates a new middleware usable by adam using the passed
// *sentry.Hub.
// For each invocation, it creates a clone of the *sentry.Hub and attaches
// meta information about the invoke.
// Afterwards, it stores the clone in the command's context under 'sentry.hub'.
func NewMiddleware(h *sentry.Hub) bot.MiddlewareFunc {
	return func(next bot.CommandFunc) bot.CommandFunc {
		return func(s *state.State, ctx *plugin.Context) error {
			h := h.Clone()

			h.WithScope(func(s *sentry.Scope) {
				s.SetTag("err_source", "bot")
			})

			h.AddBreadcrumb(&sentry.Breadcrumb{
				Category: "bot.router",
				Message:  "routed command",
				Data: map[string]interface{}{
					"command_provider": ctx.InvokedCommand.ProviderName,
					"command_id":       ctx.InvokedCommand.Identifier,
					"message":          ctx.Message,
					"lang":             ctx.Lang,
					"invoker_id":       ctx.Author.ID,
					"message_id":       ctx.Message.ID,
					"channel_id":       ctx.ChannelID,
					"guild_id":         ctx.GuildID,
				},
			}, nil)
			ctx.Set(hubKey, h)

			return next(s, ctx)
		}
	}
}
