package memory

import (
	"github.com/diamondburned/arikawa/v2/discord"

	"github.com/mavolin/levin/pkg/confgetter"
)

var _ confgetter.Repository = new(Repository)

func (d *Repository) GuildSettings(guildID discord.GuildID) (*confgetter.GuildSettings, error) {

	settings, ok := d.GuildSettingsData[guildID]
	if !ok {
		settings = &confgetter.GuildSettings{
			Language: d.defaults.Language,
			TimeZone: d.defaults.TimeZone,
			Prefix:   d.defaults.Prefix,
		}

		d.GuildSettingsData[guildID] = settings
	}

	return settings, nil
}

func (d *Repository) UserSettings(userID discord.UserID) (*confgetter.UserSettings, error) {
	settings, ok := d.UserSettingsData[userID]
	if !ok {
		settings = &confgetter.UserSettings{
			Language: d.defaults.Language,
			TimeZone: d.defaults.TimeZone,
		}

		d.UserSettingsData[userID] = settings
	}

	return settings, nil
}
