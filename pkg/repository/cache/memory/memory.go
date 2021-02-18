// Package memory provides an in-memory cache.
package memory

import (
	"sync"

	"github.com/diamondburned/arikawa/v2/discord"

	"github.com/mavolin/levin/pkg/confgetter"
	"github.com/mavolin/levin/pkg/repository/cache"
	"github.com/mavolin/levin/pkg/repository/clone"
)

type Cache struct {
	guildSettings      map[discord.GuildID]*confgetter.GuildSettings
	guildSettingsMutex sync.RWMutex

	userSettings      map[discord.UserID]*confgetter.UserSettings
	userSettingsMutex sync.RWMutex
}

var _ cache.Cache = new(Cache)

// New creates a new in-memory cache.
func New() *Cache {
	return &Cache{
		guildSettings: make(map[discord.GuildID]*confgetter.GuildSettings),
		userSettings:  make(map[discord.UserID]*confgetter.UserSettings),
	}
}

func (c *Cache) GuildSettings(guildID discord.GuildID) *confgetter.GuildSettings {
	c.guildSettingsMutex.RLock()
	defer c.guildSettingsMutex.RUnlock()

	return clone.GuildSettings(c.guildSettings[guildID])
}

func (c *Cache) SetGuildSettings(guildID discord.GuildID, s *confgetter.GuildSettings) {
	c.guildSettingsMutex.Lock()
	defer c.guildSettingsMutex.Unlock()

	c.guildSettings[guildID] = clone.GuildSettings(s)
}

func (c *Cache) UserSettings(userID discord.UserID) *confgetter.UserSettings {
	c.userSettingsMutex.RLock()
	defer c.userSettingsMutex.RUnlock()

	return clone.UserSettings(c.userSettings[userID])
}

func (c *Cache) SetUserSettings(userID discord.UserID, s *confgetter.UserSettings) {
	c.userSettingsMutex.Lock()
	defer c.userSettingsMutex.Unlock()

	c.userSettings[userID] = clone.UserSettings(s)
}
