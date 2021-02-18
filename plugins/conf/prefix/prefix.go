// Package prefix provides the prefix command
package prefix

import (
	"time"

	"github.com/diamondburned/arikawa/v2/discord"
	"github.com/mavolin/adam/pkg/impl/arg"
	"github.com/mavolin/adam/pkg/impl/command"
	"github.com/mavolin/adam/pkg/impl/throttler"
	"github.com/mavolin/adam/pkg/plugin"
	"github.com/mavolin/disstate/v3/pkg/state"

	"github.com/mavolin/levin/pkg/botmaster"
)

// Prefix is the prefix command used to display and edit the bot's prefix.
type Prefix struct {
	command.LocalizedMeta
	repo Repository
}

type Repository interface {
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
			Args: &arg.LocalizedCommaConfig{
				Optional: []arg.LocalizedOptionalArg{
					{
						Name: argNewPrefixName,
						Type: arg.Text{
							MinLength: 1,
							MaxLength: 10,
						},
					},
				},
				Flags: []arg.LocalizedFlag{
					{
						Name:        "remove",
						Aliases:     []string{"rm"},
						Type:        arg.Switch,
						Description: flagRemoveDescription,
					},
				},
			},
			Hidden:         false,
			ChannelTypes:   plugin.GuildChannels,
			BotPermissions: discord.PermissionSendMessages,
			Throttler:      throttler.PerGuild(7, time.Hour),
		},
		repo: r,
	}
}

func (p *Prefix) Invoke(s *state.State, ctx *plugin.Context) (interface{}, error) {
	newPrefix := ctx.Args.String(0)
	remove := ctx.Flags.Bool("remove")

	if remove {
		if err := botmaster.Restriction(s, ctx); err != nil {
			return nil, err
		}

		return nil, p.repo.SetPrefix(ctx.GuildID, "")
	} else if len(newPrefix) == 0 {
		if len(ctx.Prefixes) == 0 {
			return responseNoPrefix, nil
		}

		return responseCurrentPrefix.WithPlaceholders(responseCurrentPrefixPlaceholders{
			Prefix: ctx.Prefixes[0],
		}), nil
	}

	if err := p.repo.SetPrefix(ctx.GuildID, newPrefix); err != nil {
		return nil, err
	}

	return responsePrefixChanged.
		WithPlaceholders(responsePrefixChangedPlaceholders{
			NewPrefix: newPrefix,
		}), nil
}
