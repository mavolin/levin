package main

import (
	"flag"
	"os"
	"os/signal"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/mavolin/adam/pkg/bot"
	"github.com/mavolin/adam/pkg/errors"
	"github.com/mavolin/adam/pkg/impl/command/help"
	"github.com/mavolin/disstate/v3/pkg/state"
	i18nimpl "github.com/nicksnyder/go-i18n/v2/i18n"
	"go.uber.org/zap"
	"golang.org/x/text/language"

	"github.com/mavolin/levin/internal/config"
	"github.com/mavolin/levin/internal/errhandler"
	"github.com/mavolin/levin/internal/i18nwrapper"
	sentryadam "github.com/mavolin/levin/internal/sentry"
	"github.com/mavolin/levin/internal/zaplog"
	"github.com/mavolin/levin/plugins/conf"
)

var (
	debug = flag.Bool("debug", false,
		"Sets the log-level to debug and uses human-readable logs. Additionally, it disables sentry error capturing.")
	configPath       = flag.String("config", "", "A custom path to the configuration file.")
	translationsPath = flag.String("translations", "",
		"A path to a directory containing additional translation files.")
)

var log *zap.SugaredLogger

func init() {
	flag.Parse()

	zaplog.Init(*debug)
	log = zap.S().Named("startup")
	errors.Log = errhandler.CommandError()

	log.With("custom_path", *configPath).
		Info("reading config")

	if err := config.Load(*configPath); err != nil {
		log.With("err", err).
			Fatal("unable to load config")
	}

	if !(*debug) {
		if err := sentryadam.Init(); err != nil {
			log.With("err", err).
				Fatal("unable to initialize sentry")
		}
	} else {
		log.Info("debug mode: disabling sentry capturing")
	}
}

func main() {
	defer zap.S().Sync() //nolint:errcheck
	defer sentry.Flush(3 * time.Second)

	bundle := i18nimpl.NewBundle(language.English)
	err := i18nwrapper.Load(bundle, *translationsPath)
	if err != nil {
		log.With("err", err).
			Fatal("unable to load translation files")
	}

	repo, err := newRepository()
	if err != nil {
		log.With("err", err).
			Fatal("unable to initialize repository")
	}

	b, err := bot.New(bot.Options{
		Token:               config.C.Token,
		SettingsProvider:    conf.NewSettingsProvider(bundle, repo),
		Owners:              config.C.Owners,
		EditAge:             config.C.EditAge,
		AllowBot:            config.C.AllowBot,
		GatewayErrorHandler: errhandler.Gateway(zap.S(), sentry.CurrentHub()),
		StateErrorHandler:   errhandler.StateError(zap.S(), sentry.CurrentHub()),
		StatePanicHandler:   errhandler.StatePanic(zap.S(), sentry.CurrentHub()),
	})
	if err != nil {
		log.With("err", err).
			Fatal("unable to create bot")
	}

	addMiddlewares(b)
	addPlugins(b)

	log.Info("starting bot")

	b.State.MustAddHandlerOnce(func(_ *state.State, e *state.ReadyEvent) {
		log.Infof("serving as %s#%s", e.User.Username, e.User.Discriminator)
	})

	if err := b.Open(); err != nil {
		log.With("err", err).
			Fatal(err)
	}

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)

	<-sig

	log.Info("received SIGINT, exiting")

	if err := b.Close(); err != nil {
		log.With("err", err).
			Error("unable to close bot")
	}
}

func addMiddlewares(b *bot.Bot) {
	smw := sentryadam.NewMiddlewares(sentry.CurrentHub())
	b.MessageCreateMiddlewares = append(b.MessageCreateMiddlewares, smw.MessageCreateMiddleware)
	b.MessageUpdateMiddlewares = append(b.MessageUpdateMiddlewares, smw.MessageUpdateMiddleware)
	b.MustAddMiddleware(smw.BotMiddleware)

	b.MustAddMiddleware(zaplog.NewMiddlewares(zap.S()))
}

func addPlugins(b *bot.Bot) {
	b.AddCommand(help.New(help.Options{}))
}
