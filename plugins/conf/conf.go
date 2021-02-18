// Package conf provides the configuration module.
package conf

import (
	"github.com/mavolin/adam/pkg/impl/module"
	"github.com/mavolin/adam/pkg/plugin"

	"github.com/mavolin/levin/pkg/confgetter"
	"github.com/mavolin/levin/plugins/conf/prefix"
)

type Repository interface {
	confgetter.Repository
	prefix.Repository
}

// New creates a new configuration module.
func New(r Repository) plugin.Module {
	m := module.New(module.LocalizedMeta{
		Name:             "conf",
		ShortDescription: shortDescription,
	})

	m.AddCommand(prefix.New(r))

	return m
}
