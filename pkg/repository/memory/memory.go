// Package memory provides an in-memory database.
package memory

import (
	"github.com/diamondburned/arikawa/v2/discord"

	"github.com/mavolin/levin/plugins/conf"
)

type (
	Repository struct {
		GuildSettings map[discord.GuildID]*GuildSettings
		UserSettings  map[discord.UserID]*UserSettings

		defaults Defaults
	}

	Defaults struct {
		Prefixes []string
		Language string
	}

	GuildSettings struct {
		Prefixes []string
		Language string
	}

	UserSettings struct {
		Language string
	}
)

var _ conf.Repository = new(Repository)

// New creates a new in-memory database.
func New(d Defaults) *Repository {
	return &Repository{defaults: d}
}

func (d *Repository) Prefixes(guildID discord.GuildID) ([]string, error) {
	return d.guildSettings(guildID).Prefixes, nil
}

func (d *Repository) GuildLanguage(guildID discord.GuildID) (string, error) {
	return d.guildSettings(guildID).Language, nil
}

func (d *Repository) guildSettings(guildID discord.GuildID) *GuildSettings {
	settings, ok := d.GuildSettings[guildID]
	if !ok {
		settings = &GuildSettings{Language: d.defaults.Language}

		settings.Prefixes = make([]string, len(d.defaults.Prefixes))
		copy(settings.Prefixes, d.defaults.Prefixes)

		d.GuildSettings[guildID] = settings
	}
	return settings
}

func (d *Repository) UserLanguage(userID discord.UserID) (string, error) {
	return d.userSettings(userID).Language, nil
}

func (d *Repository) userSettings(userID discord.UserID) *UserSettings {
	settings, ok := d.UserSettings[userID]
	if !ok {
		settings = &UserSettings{Language: d.defaults.Language}

		d.UserSettings[userID] = settings
	}
	return settings
}
