package main

import (
	"github.com/mavolin/levin/internal/config"
	"github.com/mavolin/levin/pkg/repository"
	"github.com/mavolin/levin/pkg/repository/memory"
	"github.com/mavolin/levin/pkg/repository/mongo"
	"github.com/mavolin/levin/plugins/conf"
)

type Repository interface {
	conf.Repository
}

func newRepository() (Repository, error) {
	defaults := &repository.Defaults{
		Prefix:   config.C.DefaultPrefix,
		Language: config.C.DefaultLanguage.String(),
		TimeZone: config.C.DefaultTimeZone,
	}

	if len(config.C.Mongo.URI) == 0 {
		return memory.New(defaults), nil
	}

	return mongo.New(config.C.Mongo.URI, config.C.Mongo.DatabaseName, defaults)
}
