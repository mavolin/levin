// Package mongo provides a mongodb-based repository.
package mongo

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"

	"github.com/mavolin/levin/pkg/repository"
	"github.com/mavolin/levin/pkg/repository/cache"
	memcache "github.com/mavolin/levin/pkg/repository/cache/memory"
)

// Repository is a mongodb-based repository
type Repository struct {
	db *mongo.Database

	guildSettings *mongo.Collection
	userSettings  *mongo.Collection

	cache cache.Cache

	defaults *repository.Defaults
}

// New creates a new mongodb repository.
// It uses the passed uri to connect to the database with the passed name.
// If a new entry is created, the passed repository.Defaults will be used.
func New(uri string, dbName string, d *repository.Defaults) (*Repository, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	m, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}

	d.FillZeros()

	db := m.Database(dbName)

	return &Repository{
		db:            db,
		guildSettings: db.Collection("guild_settings"),
		userSettings:  db.Collection("user_settings"),
		cache:         memcache.New(),
		defaults:      d,
	}, nil
}

// Ping attempts to ping the database.
func (r *Repository) Ping() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	return r.db.Client().Ping(ctx, readpref.Primary())
}
