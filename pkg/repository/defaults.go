// Package repository provides sub-packages with repository implementations.
package repository

import (
	"time"

	"golang.org/x/text/language"
)

// Defaults provides default values, used if a new entity is added to a
// repository.
type Defaults struct {
	// Prefix is the default prefix used for new guilds.
	Prefix string
	// Language is the default language used for new guilds and direct
	// messages.
	Language language.Tag
	// TimeZone is the default *time.Location used for new guilds and direct
	// messages.
	TimeZone *time.Location
}
