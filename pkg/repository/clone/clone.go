// Package clone provides cloning functions for types used by Repositories.
package clone

import (
	"github.com/diamondburned/arikawa/v2/discord"

	"github.com/mavolin/levin/pkg/confgetter"
)

// UserSettings creates a clone of the passed *conf.UserSettings.
func UserSettings(s *confgetter.UserSettings) *confgetter.UserSettings {
	if s == nil {
		return nil
	}

	cl := *s
	return &cl
}

// GuildSettings creates a clone of the passed *conf.GuildSettings.
func GuildSettings(s *confgetter.GuildSettings) *confgetter.GuildSettings {
	if s == nil {
		return nil
	}

	cl := *s
	cl.BotMasterUserIDs = make([]discord.UserID, len(s.BotMasterUserIDs))
	cl.BotMasterRoleIDs = make([]discord.RoleID, len(s.BotMasterRoleIDs))

	copy(cl.BotMasterUserIDs, s.BotMasterUserIDs)
	copy(cl.BotMasterRoleIDs, s.BotMasterRoleIDs)

	return &cl
}
