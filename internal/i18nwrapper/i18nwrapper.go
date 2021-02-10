// Package i18nwrapper provides a wrappers around Nick Snyder's go-i18n
// library.
package i18nwrapper

import (
	"github.com/mavolin/adam/pkg/i18n"
	i18nimpl "github.com/nicksnyder/go-i18n/v2/i18n"
	"go.uber.org/zap"
)

func log() *zap.SugaredLogger { return zap.S().Named("startup") }

// FuncForBundle returns a i18n.Func that localizes to the passed language,
// using the passed *i18nimpl.Bundle.
func FuncForBundle(b *i18nimpl.Bundle, lang string) i18n.Func {
	l := i18nimpl.NewLocalizer(b, lang)

	return func(term i18n.Term, placeholders map[string]interface{}, plural interface{}) (string, error) {
		return l.Localize(&i18nimpl.LocalizeConfig{
			MessageID:    string(term),
			TemplateData: placeholders,
			PluralCount:  plural,
		})
	}
}
