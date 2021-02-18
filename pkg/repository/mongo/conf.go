package mongo

import (
	"context"

	"github.com/diamondburned/arikawa/v2/discord"
	"github.com/mavolin/adam/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/mavolin/levin/pkg/confgetter"
	"github.com/mavolin/levin/pkg/repository"
)

var _ confgetter.Repository = new(Repository)

// =============================================================================
// Types
// =====================================================================================

// ================================ guildSettings ================================

type guildSettings struct {
	GuildID          discord.GuildID  `bson:"guild_id"`
	Prefix           string           `bson:"prefix,omitempty"`
	Language         languageTag      `bson:"language,omitempty"`
	TimeZone         *location        `bson:"time_zone"`
	BotMasterUserIDs []discord.UserID `bson:"bot_master_user_ids"`
	BotMasterRoleIDs []discord.RoleID `bson:"bot_master_role_ids"`
}

func newDefaultGuildSettings(guildID discord.GuildID, d *repository.Defaults) *guildSettings {
	return &guildSettings{
		GuildID:  guildID,
		Prefix:   d.Prefix,
		Language: languageTag(d.Language),
		TimeZone: (*location)(d.TimeZone),
	}
}

func (s *guildSettings) toConfType() *confgetter.GuildSettings {
	return &confgetter.GuildSettings{
		Prefix:           s.Prefix,
		Language:         s.Language.baseType(),
		TimeZone:         s.TimeZone.baseType(),
		BotMasterUserIDs: s.BotMasterUserIDs,
		BotMasterRoleIDs: s.BotMasterRoleIDs,
	}
}

// ================================ userSettings ================================

type userSettings struct {
	UserID   discord.UserID `bson:"user_id"`
	Language languageTag    `bson:"language,omitempty"`
	TimeZone *location      `bson:"time_zone"`
}

func newDefaultUserSettings(userID discord.UserID, d *repository.Defaults) *userSettings {
	return &userSettings{
		UserID:   userID,
		Language: languageTag(d.Language),
		TimeZone: (*location)(d.TimeZone),
	}
}

func (s *userSettings) toConfType() *confgetter.UserSettings {
	return &confgetter.UserSettings{
		Language: s.Language.baseType(),
		TimeZone: s.TimeZone.baseType(),
	}
}

// =============================================================================
// Methods
// =====================================================================================

func (r *Repository) GuildSettings(guildID discord.GuildID) (*confgetter.GuildSettings, error) {
	if !guildID.IsValid() {
		return nil, errors.NewWithStack("repository: invalid guild id")
	}

	if s := r.cache.GuildSettings(guildID); s != nil {
		return s, nil
	}

	var s *guildSettings

	err := r.db.Client().UseSession(context.Background(), func(ctx mongo.SessionContext) error {
		res := r.guildSettings.FindOne(ctx, bson.M{"guild_id": guildID})
		if errors.Is(res.Err(), mongo.ErrNoDocuments) { // new guild?
			s = newDefaultGuildSettings(guildID, r.defaults)

			_, err := r.guildSettings.InsertOne(ctx, s)
			return errors.WithStack(err)
		} else if res.Err() != nil {
			return errors.WithStack(res.Err())
		}

		return errors.WithStack(res.Decode(&s))
	})
	if err != nil {
		return nil, errors.WithStack(err)
	}

	sconf := s.toConfType()

	r.cache.SetGuildSettings(guildID, sconf)
	return sconf, nil
}

func (r *Repository) SetPrefix(guildID discord.GuildID, prefix string) error {
	if !guildID.IsValid() {
		return errors.NewWithStack("repository: invalid guild id")
	}

	if s := r.cache.GuildSettings(guildID); s != nil {
		s.Prefix = prefix
		r.cache.SetGuildSettings(guildID, s)
	}

	err := r.db.Client().UseSession(context.Background(), func(ctx mongo.SessionContext) error {
		res, err := r.guildSettings.UpdateOne(context.Background(), bson.M{"guild_id": guildID},
			bson.M{"prefix": prefix})
		if err != nil {
			return errors.WithStack(err)
		}

		if res.MatchedCount > 0 {
			return nil
		}

		// new guild
		s := newDefaultGuildSettings(guildID, r.defaults)
		s.Prefix = prefix

		_, err = r.guildSettings.InsertOne(ctx, s)
		return errors.WithStack(err)
	})

	return errors.WithStack(err)
}

func (r *Repository) UserSettings(userID discord.UserID) (*confgetter.UserSettings, error) {
	if !userID.IsValid() {
		return nil, errors.NewWithStack("repository: invalid user id")
	}

	if s := r.cache.UserSettings(userID); s != nil {
		return s, nil
	}

	var s *userSettings

	err := r.db.Client().UseSession(context.Background(), func(ctx mongo.SessionContext) error {
		res := r.userSettings.FindOne(ctx, bson.M{"user_id": userID})
		if errors.Is(res.Err(), mongo.ErrNoDocuments) { // new user?
			s = newDefaultUserSettings(userID, r.defaults)

			_, err := r.userSettings.InsertOne(ctx, s)
			return errors.WithStack(err)
		} else if res.Err() != nil {
			return errors.WithStack(res.Err())
		}

		return errors.WithStack(res.Decode(&s))
	})
	if err != nil {
		return nil, errors.WithStack(err)
	}

	sconf := s.toConfType()

	r.cache.SetUserSettings(userID, sconf)
	return sconf, nil
}
