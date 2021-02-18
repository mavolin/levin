// Package confgetter provides helper functions to access the guild and user
// settings stored in the *plugin.Context.
package confgetter

import (
	"time"

	"github.com/diamondburned/arikawa/v2/discord"
	"github.com/mavolin/adam/pkg/impl/arg"
	"github.com/mavolin/adam/pkg/plugin"
	"golang.org/x/text/language"
)

// =============================================================================
// Always Available
// =====================================================================================

// Language returns the language of the guild or user.
// This information is available for both guild and dm invokes.
func Language(ctx *plugin.Context) language.Tag {
	if tag := ctx.Get(languageTagKey); tag != nil {
		if tag, ok := tag.(language.Tag); ok {
			return tag
		}
	}

	return language.Und
}

// TimeZone returns the time zone of the guild or user.
// This information is available for both guild and dm invokes.
func TimeZone(ctx *plugin.Context) *time.Location {
	if tz := ctx.Get(arg.LocationKey); tz != nil {
		if tz, ok := tz.(*time.Location); ok && tz != nil {
			return tz
		}
	}

	return time.UTC
}

// =============================================================================
// Guild-Specific
// =====================================================================================

// BotMasterUserIDs returns the discord.UserIDs of the users who were granted
// bot master status in the guild.
func BotMasterUserIDs(ctx *plugin.Context) []discord.UserID {
	if ids := ctx.Get(botMasterUserIDsKey); ids != nil {
		// zero-value is nil, which is fallback anyway
		ids, _ := ids.([]discord.UserID)
		return ids
	}

	return nil
}

// BotMasterRoleIDs returns the discord.RoleIDs of the roles who were granted
// bot master status in the guild.
func BotMasterRoleIDs(ctx *plugin.Context) []discord.RoleID {
	if ids := ctx.Get(botMasterRoleIDsKey); ids != nil {
		// zero-value is nil, which is fallback anyway
		ids, _ := ids.([]discord.RoleID)
		return ids
	}

	return nil
}
