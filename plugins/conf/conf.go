// Package conf provides the configuration module.
package conf

import (
	"github.com/mavolin/adam/pkg/impl/module"

	"github.com/mavolin/levin/pkg/confgetter"
)

// Configuration is the configuration module
type Configuration struct {
	*module.Module
	repo confgetter.Repository
}

// New creates a new configuration module.
func New(r confgetter.Repository) *Configuration {
	return &Configuration{
		Module: module.New(module.LocalizedMeta{
			Name:             "conf",
			ShortDescription: shortDescription,
		}),
		repo: r,
	}
}
