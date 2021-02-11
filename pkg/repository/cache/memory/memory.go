// Package memory provides an in-memory cache.
package memory

import (
	"sync"

	"github.com/diamondburned/arikawa/v2/discord"

	cache2 "github.com/mavolin/levin/pkg/repository/cache"
)

type Cache struct {
	guildSettings      map[discord.GuildID]*cache2.GuildSettings
	guildSettingsMutex sync.RWMutex

	userSettings      map[discord.UserID]*cache2.UserSettings
	userSettingsMutex sync.RWMutex
}

// New creates a new in-memory cache.
func New() *Cache {
	return &Cache{
		guildSettings: make(map[discord.GuildID]*cache2.GuildSettings),
		userSettings:  make(map[discord.UserID]*cache2.UserSettings),
	}
}

func (c *Cache) GuildSettings(guildID discord.GuildID) *cache2.GuildSettings {
	c.guildSettingsMutex.RLock()
	defer c.guildSettingsMutex.RUnlock()

	return c.guildSettings[guildID].Clone()
}

func (c *Cache) SetGuildSettings(guildID discord.GuildID, s *cache2.GuildSettings) {
	c.guildSettingsMutex.Lock()
	defer c.guildSettingsMutex.Unlock()

	c.guildSettings[guildID] = s.Clone()
}

func (c *Cache) UserSettings(userID discord.UserID) *cache2.UserSettings {
	c.userSettingsMutex.RLock()
	defer c.userSettingsMutex.RUnlock()

	return c.userSettings[userID].Clone()
}

func (c *Cache) SetUserSettings(userID discord.UserID, s *cache2.UserSettings) {
	c.userSettingsMutex.Lock()
	defer c.userSettingsMutex.Unlock()

	c.userSettings[userID] = s.Clone()
}
