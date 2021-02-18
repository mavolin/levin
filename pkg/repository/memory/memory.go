// Package memory provides an in-memory database.
package memory

import (
	"sync"

	"github.com/diamondburned/arikawa/v2/discord"

	"github.com/mavolin/levin/pkg/confgetter"
	"github.com/mavolin/levin/pkg/repository"
)

type Repository struct {
	GuildSettingsData  map[discord.GuildID]*confgetter.GuildSettings
	guildSettingsMutex sync.Mutex

	UserSettingsData  map[discord.UserID]*confgetter.UserSettings
	userSettingsMutex sync.Mutex

	defaults *repository.Defaults
}

// New creates a new in-memory database.
func New(d *repository.Defaults) *Repository {
	return &Repository{defaults: d}
}
