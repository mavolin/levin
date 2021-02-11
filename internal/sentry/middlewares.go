package sentry

import (
	"context"

	"github.com/getsentry/sentry-go"
	"github.com/mavolin/adam/pkg/bot"
	"github.com/mavolin/adam/pkg/plugin"
	"github.com/mavolin/disstate/v3/pkg/state"
)

const (
	hubKey = "sentry.hub"

	messageSpanKey     = "sentry.span.message"
	routeSpanKey       = messageSpanKey + ".route"
	middlewaresSpanKey = messageSpanKey + ".middlewares"
)

// Middlewares is a data struct that contains the middlewares sentry provides.
type Middlewares struct {
	MessageCreateMiddleware func(*state.State, *state.MessageCreateEvent)
	MessageUpdateMiddleware func(*state.State, *state.MessageUpdateEvent)
	Middleware              bot.MiddlewareFunc
	PostMiddleware          bot.MiddlewareFunc
}

// NewMiddlewares returns a struct containing all Middlewares relevant to
// sentry.
// Callers must ensure that either all, or no middlewares are added, to ensure
// proper functionality.
func NewMiddlewares(h *sentry.Hub) Middlewares {
	return Middlewares{
		MessageCreateMiddleware: routeCreateMiddleware(h),
		MessageUpdateMiddleware: routeUpdateMiddleware(h),
		Middleware:              middlewaresMiddleware(),
		PostMiddleware:          execMiddleware(),
	}
}

func routeCreateMiddleware(h *sentry.Hub) func(*state.State, *state.MessageCreateEvent) {
	return func(_ *state.State, e *state.MessageCreateEvent) {
		h := h.Clone()
		e.Set(hubKey, h)

		ctx := context.WithValue(context.Background(), sentry.HubContextKey, h)
		msgSpan := sentry.StartSpan(ctx, "message_receive")

		e.Set(messageSpanKey, msgSpan)
		e.Set(routeSpanKey, msgSpan.StartChild("route"))
	}
}

func routeUpdateMiddleware(h *sentry.Hub) func(*state.State, *state.MessageUpdateEvent) {
	return func(_ *state.State, e *state.MessageUpdateEvent) {
		h := h.Clone()
		e.Set(hubKey, h)

		ctx := context.WithValue(context.Background(), sentry.HubContextKey, h)
		span := sentry.StartSpan(ctx, messageSpanKey)

		e.Set(routeSpanKey, span)
	}
}

func middlewaresMiddleware() bot.MiddlewareFunc {
	return func(next bot.CommandFunc) bot.CommandFunc {
		return func(s *state.State, ctx *plugin.Context) error {
			h := GetHub(ctx)

			h.Scope().SetTransaction(ctx.InvokedCommand.ProviderName + "/" + string(ctx.InvokedCommand.Identifier))
			h.Scope().SetTags(map[string]string{
				"err_source":      "bot",
				"plugin_provider": ctx.InvokedCommand.ProviderName,
				"command_id":      string(ctx.InvokedCommand.Identifier),
				"lang":            ctx.Lang,
			})
			h.Scope().SetExtra("message", ctx.Message)

			if routeSpan := ctx.Get(routeSpanKey); routeSpan != nil {
				if routeSpan, ok := routeSpan.(*sentry.Span); ok && routeSpan != nil {
					routeSpan.Finish()
				}
			}

			if msgSpan := ctx.Get(messageSpanKey); msgSpan != nil {
				if msgSpan, ok := msgSpan.(*sentry.Span); ok && msgSpan != nil {
					ctx.Set(middlewaresSpanKey, msgSpan.StartChild("middlewares"))
				}
			}

			return next(s, ctx)
		}
	}
}

func execMiddleware() bot.MiddlewareFunc {
	return func(next bot.CommandFunc) bot.CommandFunc {
		return func(s *state.State, ctx *plugin.Context) error {
			if middlewaresSpan := ctx.Get(middlewaresSpanKey); middlewaresSpan != nil {
				if middlewaresSpan, ok := middlewaresSpan.(*sentry.Span); ok && middlewaresSpan != nil {
					middlewaresSpan.Finish()
				}
			}

			if msgSpan := ctx.Get(messageSpanKey); msgSpan != nil {
				if msgSpan, ok := msgSpan.(*sentry.Span); ok && msgSpan != nil {
					defer msgSpan.Finish()

					execSpan := msgSpan.StartChild("exec")
					defer execSpan.Finish()
				}
			}

			return next(s, ctx)
		}
	}
}
