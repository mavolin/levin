// Package cache provides a caching layer for mongo.
package cache

import (
	"time"

	"github.com/diamondburned/arikawa/v2/discord"
)

// Cache is the cache abstraction.
type Cache interface {
	// GuildSettings returns the settings of the guild with the passed id.
	// If the settings are not cached, it returns nil.
	GuildSettings(guildID discord.GuildID) *GuildSettings
	// SetGuildSettings updates the settings of the guild with the passed id,
	// with the passed settings.
	SetGuildSettings(guildID discord.GuildID, s *GuildSettings)
	// UserSettings returns the settings of the user with the passed id.
	// If the settings are not cached, it returns nil.
	UserSettings(userID discord.UserID) *UserSettings
	// SetUserSettings updates the settings of the user with the passed id,
	// with the passed settings.
	SetUserSettings(userID discord.UserID, s *UserSettings)
}

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

func (s *GuildSettings) Clone() *GuildSettings {
	if s == nil {
		return nil
	}

	cp := &GuildSettings{
		Language: s.Language,
		TimeZone: s.TimeZone,
	}

	cp.Prefixes = make([]string, len(s.Prefixes))
	copy(cp.Prefixes, s.Prefixes)

	return cp
}

func (s *UserSettings) Clone() *UserSettings {
	if s == nil {
		return nil
	}

	return &UserSettings{
		Language: s.Language,
		TimeZone: s.TimeZone,
	}
}
