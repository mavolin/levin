package i18nwrapper

import (
	"bytes"
	"encoding/json"
	"io/fs"
	"os"
	"regexp"

	i18nimpl "github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"

	"github.com/mavolin/levin/assets"
)

// Load loads the embedded translation files.
// Optionally, if customPath isn't empty, it will load additional translations
// found int the passed path.
// All non-translation files will be skipped.
func Load(b *i18nimpl.Bundle, customPath string) error {
	err := loadEmbeddedTranslations(b)
	if err != nil {
		return err
	}

	if len(customPath) > 0 {
		return loadCustomTranslations(b, customPath)
	}

	return nil
}

type (
	translation struct {
		Term string `json:"term"`
		// Definition is either a string or definition
		Definition json.RawMessage `json:"definition"`
	}

	// definition is the type used, if there are multiple definitions.
	definition struct {
		Zero  string `json:"zero"`
		One   string `json:"one"`
		Two   string `json:"two"`
		Few   string `json:"few"`
		Many  string `json:"many"`
		Other string `json:"other"`
	}
)

func (d *definition) isEmpty() bool {
	return len(d.Zero) == 0 && len(d.One) == 0 && len(d.Two) == 0 &&
		len(d.Few) == 0 && len(d.Many) == 0 && len(d.Other) == 0
}

var (
	defaultFileRegexp = regexp.MustCompile(`^(?P<lang>.+?)(?:_(?:adam|levin))?\.json$`)
	customFileRegexp  = regexp.MustCompile(`^(?P<lang>.+?)\.json$`)
)

func loadEmbeddedTranslations(b *i18nimpl.Bundle) error {
	dir, err := assets.Translations.ReadDir("translations")
	if err != nil {
		return err
	}

	for _, f := range dir {
		if f.IsDir() {
			continue
		}

		matches := defaultFileRegexp.FindStringSubmatch(f.Name())
		if len(matches) < 2 {
			log().With("file_name", f.Name()).
				Warn("found non-translation file in embedded translations, skipping")
		}

		tag, err := language.Parse(matches[1])
		if err != nil {
			log().With("lang", matches[1]).
				Warn("embedded translations contain a translation file for invalid language, skipping")
		}

		f, err := assets.Translations.Open("translations/" + f.Name())
		if err != nil {
			return err
		}

		if err = loadTranslation(b, tag, f); err != nil {
			return err
		}
	}

	return err
}

func loadCustomTranslations(b *i18nimpl.Bundle, customPath string) error {
	dir, err := os.Open(customPath)
	if err != nil {
		return err
	}

	files, err := dir.Readdir(0)
	if err != nil {
		return err
	}

	for _, f := range files {
		if f.IsDir() {
			continue
		}

		matches := customFileRegexp.FindStringSubmatch(f.Name())
		if len(matches) < 2 {
			continue
		}

		tag, err := language.Parse(matches[1])
		if err != nil {
			continue
		}

		f, err := os.Open(customPath + f.Name())
		if err != nil {
			return err
		}

		if err = loadTranslation(b, tag, f); err != nil {
			return err
		}
	}

	return nil
}

func loadTranslation(b *i18nimpl.Bundle, tag language.Tag, f fs.File) error {
	var messages []translation

	if err := json.NewDecoder(f).Decode(&messages); err != nil {
		return err
	}

	for _, m := range messages {
		var def definition

		if bytes.HasPrefix(m.Definition, []byte(`"`)) { // just other
			if err := json.Unmarshal(m.Definition, &def.Other); err != nil {
				return err
			}
		} else { // definition object
			if err := json.Unmarshal(m.Definition, &def); err != nil {
				return err
			}
		}

		if len(m.Term) == 0 || def.isEmpty() {
			continue
		}

		err := b.AddMessages(tag, &i18nimpl.Message{
			ID:    m.Term,
			Zero:  def.Zero,
			One:   def.One,
			Two:   def.Two,
			Few:   def.Few,
			Many:  def.Many,
			Other: def.Other,
		})
		if err != nil {
			return err
		}
	}

	return nil
}
