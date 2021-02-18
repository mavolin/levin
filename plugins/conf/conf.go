// Package conf provides the configuration module.
package conf

import (
	"github.com/mavolin/adam/pkg/impl/module"
	"github.com/mavolin/adam/pkg/plugin"
	i18nimpl "github.com/nicksnyder/go-i18n/v2/i18n"

	"github.com/mavolin/levin/pkg/confgetter"
	"github.com/mavolin/levin/plugins/conf/language"
	"github.com/mavolin/levin/plugins/conf/prefix"
)

type Repository interface {
	confgetter.Repository

	prefix.Repository
	language.Repository
}

// New creates a new configuration module.
func New(r Repository, bundle *i18nimpl.Bundle) plugin.Module {
	m := module.New(module.LocalizedMeta{
		Name:             "conf",
		ShortDescription: shortDescription,
	})

	m.AddCommand(prefix.New(r))
	m.AddCommand(language.New(bundle, r))

	return m
}
