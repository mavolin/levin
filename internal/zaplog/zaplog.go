// Package zaplog provides zap log wrappers for adam.
package zaplog

import (
	"log"

	"github.com/mavolin/adam/pkg/bot"
	"github.com/mavolin/adam/pkg/plugin"
	"github.com/mavolin/disstate/v3/pkg/state"
	jww "github.com/spf13/jwalterweatherman"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const loggerKey = "logger"

// Init initializes the global zap logger.
func Init(debug bool) {
	if debug {
		l, err := zap.NewDevelopment()
		if err != nil {
			panic(err)
		}

		zap.ReplaceGlobals(l)
	} else {
		l, err := zap.NewProduction()
		if err != nil {
			panic(err)
		}

		zap.ReplaceGlobals(l)
	}

	jww.TRACE = mustStdLogAt(zap.L(), zapcore.DebugLevel)
	jww.DEBUG = mustStdLogAt(zap.L(), zapcore.DebugLevel)
	// viper's info logs are very verbose, using debug for this
	jww.DEBUG = mustStdLogAt(zap.L(), zapcore.DebugLevel)
	jww.WARN = mustStdLogAt(zap.L(), zapcore.WarnLevel)
	jww.ERROR = mustStdLogAt(zap.L(), zapcore.ErrorLevel)
	jww.CRITICAL = mustStdLogAt(zap.L(), zapcore.ErrorLevel)
	jww.FATAL = mustStdLogAt(zap.L(), zapcore.FatalLevel)
	jww.LOG = mustStdLogAt(zap.L(), zapcore.InfoLevel)
}

func mustStdLogAt(l *zap.Logger, lvl zapcore.Level) *log.Logger {
	stdl, err := zap.NewStdLogAt(l, lvl)
	if err != nil {
		zap.S().Named("config").Fatal(err)
	}

	return stdl
}

// Get extracts a *zap.SugaredLogger from the Context.
// This will only work, if a *zap.SugaredLogger was previously stored in the
// Context under 'logger', e.g. using NewMiddleware.
func Get(ctx *plugin.Context) *zap.SugaredLogger {
	return ctx.Get(loggerKey).(*zap.SugaredLogger)
}

// NewMiddlewares creates a new bot.MiddlewareFunc that stores the passed
// *zap.SugaredLogger under 'logger' in the command's plugin.Context.
// Additionally, it attaches some meta information to the logger.
func NewMiddlewares(l *zap.SugaredLogger) bot.MiddlewareFunc {
	l = l.Named("bot")

	return func(next bot.CommandFunc) bot.CommandFunc {
		return func(s *state.State, ctx *plugin.Context) error {
			l := l.With(
				"command_provider", ctx.InvokedCommand.ProviderName,
				"command_id", ctx.InvokedCommand.Identifier,
				"message", ctx.Content,
				"lang", ctx.Lang,
				"invoker_id", ctx.Author.ID,
				"message_id", ctx.ID,
				"channel_id", ctx.ChannelID,
				"guild_id", ctx.GuildID,
			)
			ctx.Set(loggerKey, l)

			l.Info("received invoke")

			return next(s, ctx)
		}
	}
}
