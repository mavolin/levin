// Package conf provides the configuration module.
package conf

import (
	"time"

	"github.com/diamondburned/arikawa/v2/discord"
	"github.com/mavolin/adam/pkg/bot"
	"github.com/mavolin/adam/pkg/impl/arg"
	"github.com/mavolin/adam/pkg/impl/module"
	"github.com/mavolin/adam/pkg/plugin"
	"github.com/mavolin/disstate/v3/pkg/state"
)

type Repository interface {
	// Prefix returns the prefix of the guild with the passed id.
	// If the returned prefix is empty, the guild is considered to only use
	// mentions.
	Prefix(guildID discord.GuildID) (string, error)
	// GuildLanguage returns the BCP 47 language tag of the guild with the
	// passed id.
	GuildLanguage(guildID discord.GuildID) (string, error)
	// UserLanguage returns the BCP 47 language tag of the user with the passed
	// id.
	UserLanguage(userID discord.UserID) (string, error)
	// GuildTimezone returns the IANA timezone identifier of the guild's
	// timezone.
	GuildTimezone(guildID discord.GuildID) (*time.Location, error)
	// UserTimezone returns the IANA timezone identifier of the user's
	// timezone.
	UserTimezone(userID discord.UserID) (*time.Location, error)
}

// Configuration is the configuration module
type Configuration struct {
	*module.Module
	repo Repository
}

// New creates a new configuration module.
func New(r Repository) *Configuration {
	return &Configuration{
		Module: module.New(module.LocalizedMeta{
			Name:             "conf",
			ShortDescription: shortDescription,
		}),
		repo: r,
	}
}

// Open adds a timezone middleware to the bot and disables time zone fallbacks
// in case the location could not be loaded.
func (c *Configuration) Open(b *bot.Bot) {
	arg.DefaultLocation = nil
	arg.LocationKey = "location"
	b.MustAddMiddleware(newTimezoneMiddleware(c.repo))
}

func newTimezoneMiddleware(r Repository) bot.MiddlewareFunc {
	return func(next bot.CommandFunc) bot.CommandFunc {
		return func(s *state.State, ctx *plugin.Context) error {
			var tz *time.Location
			var err error

			if ctx.GuildID == 0 {
				tz, err = r.UserTimezone(ctx.Author.ID)
			} else {
				tz, err = r.GuildTimezone(ctx.GuildID)
			}

			if err != nil {
				ctx.HandleErrorSilently(err)
				return next(s, ctx)
			}

			ctx.Set(arg.LocationKey, tz)
			return next(s, ctx)
		}
	}
}
