// Package repository provides sub-packages with repository implementations.
package repository

import "time"

// Defaults provides default values, used if a new entity is added to a
// repository.
type Defaults struct {
	// Prefixes are the default prefixes used for new guilds.
	Prefixes []string
	// Language is the default language used for new guilds and direct
	// messages.
	Language string
	// TimeZone is the default *time.Location used for new guilds and direct
	// messages.
	TimeZone *time.Location
}

// FillZeros fills all empty/invalid defaults.
func (d *Defaults) FillZeros() {
	if len(d.Language) == 0 {
		d.Language = "en"
	}

	if d.TimeZone == nil {
		d.TimeZone = time.UTC
	}
}
