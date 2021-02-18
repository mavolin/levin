package memory

import (
	"github.com/diamondburned/arikawa/v2/discord"
	"github.com/mavolin/adam/pkg/errors"

	"github.com/mavolin/levin/pkg/confgetter"
)

var _ confgetter.Repository = new(Repository)

func (d *Repository) GuildSettings(guildID discord.GuildID) (*confgetter.GuildSettings, error) {
	if !guildID.IsValid() {
		return nil, errors.NewWithStack("repository: invalid guild id")
	}

	d.guildSettingsMutex.Lock()
	defer d.guildSettingsMutex.Unlock()

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

func (d *Repository) SetPrefix(guildID discord.GuildID, prefix string) error {
	if !guildID.IsValid() {
		return errors.NewWithStack("repository: invalid guild id")
	}

	d.guildSettingsMutex.Lock()
	defer d.guildSettingsMutex.Unlock()

	settings, ok := d.GuildSettingsData[guildID]
	if !ok {
		settings = &confgetter.GuildSettings{
			Language: d.defaults.Language,
			TimeZone: d.defaults.TimeZone,
			Prefix:   prefix,
		}

		d.GuildSettingsData[guildID] = settings
	} else {
		settings.Prefix = prefix
		d.GuildSettingsData[guildID] = settings
	}

	return nil
}

func (d *Repository) UserSettings(userID discord.UserID) (*confgetter.UserSettings, error) {
	if !userID.IsValid() {
		return nil, errors.NewWithStack("repository: invalid user id")
	}

	d.userSettingsMutex.Lock()
	defer d.userSettingsMutex.Unlock()

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
