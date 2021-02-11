package memory

import (
	"time"

	"github.com/diamondburned/arikawa/v2/discord"

	"github.com/mavolin/levin/plugins/conf"
)

type (
	GuildSettings struct {
		Prefixes []string
		Language string
		TimeZone *time.Location
	}

	UserSettings struct {
		Language string
		TimeZone *time.Location
	}
)

var _ conf.Repository = new(Repository)

func (d *Repository) Prefixes(guildID discord.GuildID) ([]string, error) {
	return d.guildSettings(guildID).Prefixes, nil
}

func (d *Repository) GuildLanguage(guildID discord.GuildID) (string, error) {
	return d.guildSettings(guildID).Language, nil
}

func (d *Repository) GuildTimezone(guildID discord.GuildID) (*time.Location, error) {
	return d.guildSettings(guildID).TimeZone, nil
}

func (d *Repository) guildSettings(guildID discord.GuildID) *GuildSettings {
	settings, ok := d.GuildSettings[guildID]
	if !ok {
		settings = &GuildSettings{
			Language: d.defaults.Language,
			TimeZone: d.defaults.TimeZone,
		}

		settings.Prefixes = make([]string, len(d.defaults.Prefixes))
		copy(settings.Prefixes, d.defaults.Prefixes)

		d.GuildSettings[guildID] = settings
	}
	return settings
}

func (d *Repository) UserLanguage(userID discord.UserID) (string, error) {
	return d.userSettings(userID).Language, nil
}

func (d *Repository) UserTimezone(userID discord.UserID) (*time.Location, error) {
	return d.userSettings(userID).TimeZone, nil
}

func (d *Repository) userSettings(userID discord.UserID) *UserSettings {
	settings, ok := d.UserSettings[userID]
	if !ok {
		settings = &UserSettings{
			Language: d.defaults.Language,
			TimeZone: d.defaults.TimeZone,
		}

		d.UserSettings[userID] = settings
	}
	return settings
}
