package conf

import (
	"github.com/diamondburned/arikawa/v2/discord"
	"github.com/mavolin/adam/pkg/bot"
	"github.com/mavolin/adam/pkg/i18n"
	"github.com/mavolin/disstate/v3/pkg/state"
	i18nimpl "github.com/nicksnyder/go-i18n/v2/i18n"

	"github.com/mavolin/levin/internal/i18nwrapper"
	sentryadam "github.com/mavolin/levin/internal/sentry"
)

// NewSettingsProvider creates a new bot.SettingsProvider using the passed
// *i18nimpl.Bundle and Repository.
func NewSettingsProvider(bundle *i18nimpl.Bundle, r Repository) bot.SettingsProvider {
	return func(b *state.Base, m *discord.Message) ([]string, *i18n.Localizer) {
		var prefixes []string

		if m.GuildID > 0 {
			var err error

			prefixes, err = r.Prefixes(m.GuildID)
			if err != nil {
				prefixes = nil
				sentryadam.Get(b).CaptureException(err)
			}
		}

		var lang string
		var err error

		if m.GuildID == 0 {
			lang, err = r.UserLanguage(m.Author.ID)
		} else {
			lang, err = r.GuildLanguage(m.GuildID)
		}

		if err != nil {
			sentryadam.Get(b).CaptureException(err)
			return prefixes, i18n.NewFallbackLocalizer()
		}

		f := i18nwrapper.FuncForBundle(bundle, lang)
		return prefixes, i18n.NewLocalizer(lang, f)
	}
}
