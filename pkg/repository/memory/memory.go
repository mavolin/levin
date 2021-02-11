// Package memory provides an in-memory database.
package memory

import (
	"github.com/diamondburned/arikawa/v2/discord"

	"github.com/mavolin/levin/pkg/repository"
)

type Repository struct {
	GuildSettings map[discord.GuildID]*GuildSettings
	UserSettings  map[discord.UserID]*UserSettings

	defaults *repository.Defaults
}

// New creates a new in-memory database.
func New(d *repository.Defaults) *Repository {
	d.FillZeros()
	return &Repository{defaults: d}
}
