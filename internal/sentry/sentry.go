package sentry

import (
	"github.com/getsentry/sentry-go"

	"github.com/mavolin/levin/internal/config"
	"github.com/mavolin/levin/internal/meta"
)

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

// GetHub extracts the *sentry.Hub from the passed context.
// The context is abstracted as interface { Get(string) interface{} }, to allow
// both getting from *state.XEvents and *plugin.Contexts.
//
// If no hub is accessible using 'sentry.hub' as key, a clone of the current
// hub will be returned.
func GetHub(ctx interface{ Get(string) interface{} }) *sentry.Hub {
	if hub := ctx.Get(hubKey); hub != nil {
		if hub, ok := hub.(*sentry.Hub); ok && hub != nil {
			return hub
		}
	}

	return sentry.CurrentHub().Clone()
}
