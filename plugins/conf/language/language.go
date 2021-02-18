// Package language provides the language command.
package language

import (
	"strings"
	"time"

	"github.com/diamondburned/arikawa/v2/discord"
	"github.com/mavolin/adam/pkg/errors"
	"github.com/mavolin/adam/pkg/impl/arg"
	"github.com/mavolin/adam/pkg/impl/command"
	"github.com/mavolin/adam/pkg/impl/throttler"
	"github.com/mavolin/adam/pkg/plugin"
	"github.com/mavolin/disstate/v3/pkg/state"
	i18nimpl "github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"

	"github.com/mavolin/levin/pkg/botmaster"
)

// Language is a plugin.Command used to edit the bots language.
type Language struct {
	command.LocalizedMeta
	repo   Repository
	bundle *i18nimpl.Bundle
}

type Repository interface {
	// SetGuildLanguage sets the language of the guild with the passed id to
	// the passed language.
	SetGuildLanguage(guildID discord.GuildID, lang language.Tag) error
	// SetUserLanguage sets the language of the user with the passed id to
	// the passed language.
	SetUserLanguage(userID discord.UserID, lang language.Tag) error
}

func New(bundle *i18nimpl.Bundle, r Repository) *Language {
	return &Language{
		LocalizedMeta: command.LocalizedMeta{
			Name:             "language",
			Aliases:          []string{"lang"},
			ShortDescription: shortDescription,
			LongDescription:  longDescription,
			ExampleArgs:      exampleArgs,
			Args: &arg.LocalizedCommaConfig{
				Optional: []arg.LocalizedOptionalArg{
					{
						Name: argNewLanguageName,
						Type: arg.SimpleText,
					},
				},
			},
			Hidden:         false,
			BotPermissions: discord.PermissionSendMessages,
			Throttler:      throttler.PerGuild(7, time.Hour),
			Restrictions:   botmaster.Restriction,
		},
		repo:   r,
		bundle: bundle,
	}
}

func (p *Language) Invoke(_ *state.State, ctx *plugin.Context) (interface{}, error) {
	newLanguage := ctx.Args.String(0)

	if len(newLanguage) == 0 {
		var langListBuilder strings.Builder
		langListBuilder.Grow(2024) // max msg len

		for i, tag := range p.bundle.LanguageTags() {
			if i > 0 {
				langListBuilder.WriteString(", ")
			}

			langListBuilder.WriteRune('`')
			langListBuilder.WriteString(tag.String())
			langListBuilder.WriteRune('`')
		}

		return responseList.WithPlaceholders(responseListPlaceholders{
			Languages: langListBuilder.String(),
		}), nil
	}

	tag, err := language.Parse(newLanguage)
	if err != nil {
		errmsg := errorInvalidLanguage.WithPlaceholders(errorInvalidLanguagePlaceholders{
			Raw:            newLanguage,
			LanguageInvoke: ctx.InvokedCommand.Identifier.AsInvoke(),
		})

		return nil, plugin.NewArgumentErrorl(errmsg)
	}

	var has bool

	for _, target := range p.bundle.LanguageTags() {
		if target.String() == tag.String() {
			has = true
			break
		}
	}

	if !has {
		errmsg := errorInvalidLanguage.WithPlaceholders(errorInvalidLanguagePlaceholders{
			Raw:            newLanguage,
			LanguageInvoke: ctx.InvokedCommand.Identifier.AsInvoke(),
		})

		return nil, plugin.NewArgumentErrorl(errmsg)
	}

	if ctx.GuildID == 0 {
		err := p.repo.SetUserLanguage(ctx.Author.ID, tag)
		if err != nil {
			return nil, errors.WithStack(err)
		}
	} else {
		err := p.repo.SetGuildLanguage(ctx.GuildID, tag)
		if err != nil {
			return nil, errors.WithStack(err)
		}
	}

	return responseLanguageChanged, nil
}
