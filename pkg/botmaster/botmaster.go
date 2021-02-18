package botmaster

import (
	"github.com/diamondburned/arikawa/v2/discord"
	"github.com/mavolin/adam/pkg/impl/restriction"
	"github.com/mavolin/adam/pkg/plugin"
	"github.com/mavolin/adam/pkg/utils/permutil"
	"github.com/mavolin/disstate/v3/pkg/state"

	"github.com/mavolin/levin/pkg/confgetter"
)

// Is checks if the invoking user is a BotMaster in the guild with the passed
// context.
func Is(m *discord.Member, ctx *plugin.Context) (bool, error) {
	for _, id := range confgetter.BotMasterUserIDs(ctx) {
		if m.User.ID == id {
			return true, nil
		}
	}

	for _, targetID := range confgetter.BotMasterRoleIDs(ctx) {
		for _, id := range m.RoleIDs {
			if targetID == id {
				return true, nil
			}
		}
	}

	g, err := ctx.Guild()
	if err != nil {
		return false, err
	}

	perms := permutil.MemberPermissions(*g, *m)

	return perms.Has(discord.PermissionAdministrator), nil
}

// Restriction is the restriction used to assert that the invoking user is a
// bot master.
func Restriction(s *state.State, ctx *plugin.Context) error {
	err := restriction.ChannelTypes(plugin.GuildChannels)(s, ctx)
	if err != nil {
		return err
	}

	is, err := Is(ctx.Member, ctx)
	if err != nil {
		return err
	}

	if is {
		return nil
	}

	return plugin.NewFatalRestrictionErrorl(botMasterError)
}

var _ plugin.RestrictionFunc = Restriction
