// Package prefix provides the prefix command
package prefix

import (
	"github.com/diamondburned/arikawa/v2/discord"
	"github.com/mavolin/adam/pkg/impl/command"
	"github.com/mavolin/adam/pkg/plugin"
)

// Prefix is the prefix command used to display and edit the bot's prefix.
type Prefix struct {
	command.LocalizedMeta
	repo Repository
}

type Repository interface {
	// Prefix returns the prefix of the guild with the passed id.
	// If the returned prefix is empty, the guild is considered to only use
	// mentions.
	Prefix(guildID discord.GuildID) (string, error)
	// SetPrefix sets the prefix of the guild with the passed id to the passed
	// prefix.
	// If the prefix is empty, only mentions will be allowed.
	SetPrefix(guildID discord.GuildID, prefix string) error
}

func New(r Repository) *Prefix {
	return &Prefix{
		LocalizedMeta: command.LocalizedMeta{
			Name:             "prefix",
			Aliases:          nil,
			ShortDescription: shortDescription,
			LongDescription:  longDescription,
			ExampleArgs:      exampleArgs,
			Args:             nil,
			Hidden:           false,
			ChannelTypes:     plugin.GuildChannels,
			BotPermissions:   discord.PermissionSendMessages,
			Restrictions:     nil,
			Throttler:        nil,
		},
		repo: r,
	}
}
