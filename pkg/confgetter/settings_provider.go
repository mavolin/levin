package confgetter

import (
	"time"

	"github.com/diamondburned/arikawa/v2/discord"
	"github.com/mavolin/adam/pkg/bot"
	"github.com/mavolin/adam/pkg/i18n"
	"github.com/mavolin/adam/pkg/impl/arg"
	"github.com/mavolin/disstate/v3/pkg/state"
	i18nimpl "github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"

	"github.com/mavolin/levin/internal/i18nwrapper"
	sentryadam "github.com/mavolin/levin/internal/sentry"
)

const (
	baseKey = "settings"

	languageTagKey      = baseKey + ".language_tag"
	botMasterUserIDsKey = baseKey + ".bot_master_user_ids"
	botMasterRoleIDsKey = baseKey + ".bot_master_role_ids"
)

type (
	Repository interface {
		// GuildSettings returns the settings for the guild with the passed id.
		GuildSettings(guildID discord.GuildID) (*GuildSettings, error)
		// UserSettings returns the settings for the user with the passed id.
		UserSettings(userID discord.UserID) (*UserSettings, error)
	}

	// GuildSettings are the base settings of a guild.
	GuildSettings struct {
		// Prefix is the prefix of the guild.
		// If the Prefix is empty, the guild uses mentions only.
		Prefix string
		// Language is the language of the guild.
		Language language.Tag
		// TimeZone is the the *time.Location of the guild.
		TimeZone *time.Location
		// BotMasterUserIDs are the discord.UserIDs of the users declared bot
		// masters by the admins of the server.
		BotMasterUserIDs []discord.UserID
		// BotMasterRoleIDs are the discord.RoleIDs of the roles declared bot
		// masters by the admins of the server.
		BotMasterRoleIDs []discord.RoleID
	}

	// UserSettings are the base settings of a user.
	UserSettings struct {
		// Language is the language of the user.
		Language language.Tag
		// TimeZone is the the *time.Location of the user.
		TimeZone *time.Location
	}
)

// NewSettingsProvider creates a new bot.SettingsProvider using the passed
// Repository and *i18nimpl.Bundle.
// It will also set the arg.LocationKey and disable DefaultLocations.
func NewSettingsProvider(r Repository, bundle *i18nimpl.Bundle) bot.SettingsProvider {
	arg.DefaultLocation = nil
	arg.LocationKey = "settings.time_zone"

	return func(b *state.Base, m *discord.Message) (prefixes []string, localizer *i18n.Localizer, ok bool) {
		if m.GuildID == 0 {
			settings, err := r.UserSettings(m.Author.ID)
			if err != nil {
				sentryadam.Get(b).CaptureException(err)
				return nil, nil, false
			}

			f := i18nwrapper.FuncForBundle(bundle, settings.Language.String())
			localizer = i18n.NewLocalizer(settings.Language.String(), f)

			b.Set(languageTagKey, settings.Language)
			b.Set(arg.LocationKey, settings.TimeZone)
		} else {
			settings, err := r.GuildSettings(m.GuildID)
			if err != nil {
				sentryadam.Get(b).CaptureException(err)
				return nil, nil, false
			}

			f := i18nwrapper.FuncForBundle(bundle, settings.Language.String())
			localizer = i18n.NewLocalizer(settings.Language.String(), f)

			if len(settings.Prefix) > 0 {
				prefixes = []string{settings.Prefix}
			}

			b.Set(languageTagKey, settings.Language)
			b.Set(arg.LocationKey, settings.TimeZone)
			b.Set(botMasterUserIDsKey, settings.BotMasterUserIDs)
			b.Set(botMasterRoleIDsKey, settings.BotMasterRoleIDs)
		}

		return prefixes, localizer, true
	}
}
