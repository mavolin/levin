package mongo

import (
	"context"
	"time"

	"github.com/diamondburned/arikawa/v2/discord"
	"github.com/mavolin/adam/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/mavolin/levin/pkg/repository"
	"github.com/mavolin/levin/pkg/repository/cache"
	"github.com/mavolin/levin/plugins/conf"
)

var _ conf.Repository = new(Repository)

// =============================================================================
// Types
// =====================================================================================

// ================================ guildSettings ================================

type guildSettings struct {
	GuildID  discord.GuildID `bson:"guild_id"`
	Prefix   string          `bson:"prefix,omitempty"`
	Language string          `bson:"language,omitempty"`
	TimeZone *Location       `bson:"time_zone"`
}

func newDefaultGuildSettings(d *repository.Defaults) *guildSettings {
	return &guildSettings{
		Prefix:   d.Prefix,
		Language: d.Language,
		TimeZone: (*Location)(d.TimeZone),
	}
}

func newGuildSettingsFromCache(s *cache.GuildSettings) *guildSettings {
	return &guildSettings{
		Prefix:   s.Prefix,
		Language: s.Language,
		TimeZone: (*Location)(s.TimeZone),
	}
}

func (s *guildSettings) CacheType() *cache.GuildSettings {
	return &cache.GuildSettings{
		Prefix:   s.Prefix,
		Language: s.Language,
		TimeZone: s.TimeZone.Location(),
	}
}

// ================================ userSettings ================================

type userSettings struct {
	UserID   discord.UserID `bson:"user_id"`
	Language string         `bson:"language,omitempty"`
	TimeZone *Location      `bson:"time_zone"`
}

func newDefaultUserSettings(d *repository.Defaults) *userSettings {
	return &userSettings{
		Language: d.Language,
		TimeZone: (*Location)(d.TimeZone),
	}
}

func newUserSettingsFromCache(s *cache.UserSettings) *userSettings {
	return &userSettings{
		Language: s.Language,
		TimeZone: (*Location)(s.TimeZone),
	}
}

func (s *userSettings) CacheType() *cache.UserSettings {
	return &cache.UserSettings{
		Language: s.Language,
		TimeZone: s.TimeZone.Location(),
	}
}

// =============================================================================
// Methods
// =====================================================================================

func (r *Repository) Prefix(guildID discord.GuildID) (string, error) {
	s, err := r.getGuildSettings(guildID)
	if err != nil {
		return "", err
	}

	return s.Prefix, nil
}

func (r *Repository) GuildLanguage(guildID discord.GuildID) (string, error) {
	s, err := r.getGuildSettings(guildID)
	if err != nil {
		return "", err
	}

	return s.Language, nil
}

func (r *Repository) GuildTimeZone(guildID discord.GuildID) (*time.Location, error) {
	s, err := r.getGuildSettings(guildID)
	if err != nil {
		return nil, err
	}

	return s.TimeZone.Location(), nil
}

func (r *Repository) UserLanguage(userID discord.UserID) (string, error) {
	s, err := r.getUserSettings(userID)
	if err != nil {
		return "", err
	}

	return s.Language, nil
}

func (r *Repository) UserTimeZone(userID discord.UserID) (*time.Location, error) {
	s, err := r.getUserSettings(userID)
	if err != nil {
		return nil, err
	}

	return s.TimeZone.Location(), nil
}

func (r *Repository) getGuildSettings(guildID discord.GuildID) (*guildSettings, error) {
	if !guildID.IsValid() {
		return nil, errors.NewWithStack("invalid guild id")
	}

	if s := r.cache.GuildSettings(guildID); s != nil {
		return newGuildSettingsFromCache(s), nil
	}

	res := r.guildSettings.FindOne(context.Background(), bson.M{"guild_id": guildID})
	if errors.Is(res.Err(), mongo.ErrNoDocuments) {
		s := newDefaultGuildSettings(r.defaults)
		s.GuildID = guildID

		_, err := r.guildSettings.InsertOne(context.Background(), s)
		if err != nil {
			return nil, errors.WithStack(err)
		}
	}

	var s *guildSettings
	if err := res.Decode(&s); err != nil {
		return nil, errors.WithStack(err)
	}

	r.cache.SetGuildSettings(guildID, s.CacheType())
	return s, nil
}

func (r *Repository) getUserSettings(userID discord.UserID) (*userSettings, error) {
	if !userID.IsValid() {
		return nil, errors.NewWithStack("invalid user id")
	}

	if s := r.cache.UserSettings(userID); s != nil {
		return newUserSettingsFromCache(s), nil
	}

	res := r.userSettings.FindOne(context.Background(), bson.M{"user_id": userID})
	if errors.Is(res.Err(), mongo.ErrNoDocuments) {
		s := newDefaultUserSettings(r.defaults)
		s.UserID = userID

		_, err := r.userSettings.InsertOne(context.Background(), s)
		if err != nil {
			return nil, errors.WithStack(err)
		}
	}

	var s *userSettings
	if err := res.Decode(&s); err != nil {
		return nil, errors.WithStack(err)
	}

	r.cache.SetUserSettings(userID, s.CacheType())
	return s, nil
}
