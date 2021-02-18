// Package cache provides a caching layer for mongo.
package cache

import (
	"github.com/diamondburned/arikawa/v2/discord"

	"github.com/mavolin/levin/pkg/confgetter"
)

// Cache is the cache abstraction.
type Cache interface {
	// GuildSettings returns the settings of the guild with the passed id.
	// If the settings are not cached, it returns nil.
	GuildSettings(guildID discord.GuildID) *confgetter.GuildSettings
	// SetGuildSettings updates the settings of the guild with the passed id,
	// with the passed settings.
	SetGuildSettings(guildID discord.GuildID, s *confgetter.GuildSettings)
	// UserSettings returns the settings of the user with the passed id.
	// If the settings are not cached, it returns nil.
	UserSettings(userID discord.UserID) *confgetter.UserSettings
	// SetUserSettings updates the settings of the user with the passed id,
	// with the passed settings.
	SetUserSettings(userID discord.UserID, s *confgetter.UserSettings)
}
