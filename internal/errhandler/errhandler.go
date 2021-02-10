// Package errhandler provides the default error handlers.
package errhandler

import (
	"runtime"
	"strconv"
	"strings"

	"github.com/getsentry/sentry-go"
	"github.com/mavolin/adam/pkg/bot"
	"github.com/mavolin/adam/pkg/plugin"
	"go.uber.org/zap"

	sentryadam "github.com/mavolin/levin/internal/sentry"
	"github.com/mavolin/levin/internal/zaplog"
)

// CommandError returns the logger function used for errors.Log.
// It logs the error using the passed *zap.SugaredLogger, and extracts the
// assigned *sentry.Hub and *zap.SugaredLogger from the context using sentryadam.GetHub.
func CommandError() func(error, *plugin.Context) {
	return func(err error, ctx *plugin.Context) {
		sentryadam.GetHub(ctx).CaptureException(err)

		l := zaplog.Get(ctx).With("err", err)

		if serr, ok := err.(interface{ StackTrace() []uintptr }); ok {
			l.With("stack_trace", writtenStack(serr.StackTrace()))
		}

		l.Error("error during command execution")
	}
}

func writtenStack(callers []uintptr) string {
	frames := runtime.CallersFrames(callers)

	first := true

	var b strings.Builder
	b.Grow(1024)

	// the last frame is skipped, as it only contains runtime information
	for frame, more := frames.Next(); more; frame, more = frames.Next() {
		if first {
			b.WriteByte('\n')
			first = false
		}

		b.WriteString(frame.Function)
		b.WriteRune('\n')
		b.WriteRune('\t')
		b.WriteString(frame.File)
		b.WriteRune(':')
		b.WriteString(strconv.Itoa(frame.Line))
	}

	return b.String()
}

// Gateway returns the error handler function used for the
// bot.Options.GatwayErrorHandler.
func Gateway(l *zap.SugaredLogger, h *sentry.Hub) func(error) {
	l = l.Named("gateway")

	h = h.Clone()
	h.WithScope(func(s *sentry.Scope) {
		s.SetTag("err_source", "gateway")
	})

	return func(err error) {
		if bot.FilterGatewayError(err) {
			h.CaptureException(err)
			l.Error(err)
		}
	}
}

// StateError returns the logger function used for the
// bot.Option.StateErrorHandler.
func StateError(l *zap.SugaredLogger, h *sentry.Hub) func(error) {
	l = l.Named("state")

	h = h.Clone()
	h.WithScope(func(s *sentry.Scope) {
		s.SetTag("err_source", "state")
	})

	return func(err error) {
		h.CaptureException(err)
		l.Error(err)
	}
}

// StatePanic returns the logger function used for the
// bot.Option.StatePanicHandler.
func StatePanic(l *zap.SugaredLogger, h *sentry.Hub) func(interface{}) {
	l = l.Named("state")

	h = h.Clone()
	h.WithScope(func(s *sentry.Scope) {
		s.SetTag("err_source", "state")
	})

	return func(rec interface{}) {
		h.Recover(rec)
		l.Errorf("recovered from panic: %+v", rec)
	}
}
