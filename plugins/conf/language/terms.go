package language

import "github.com/mavolin/adam/pkg/i18n"

// =============================================================================
// Meta
// =====================================================================================

var (
	shortDescription = i18n.NewFallbackConfig(
		"plugin.conf.language.short_description",
		"🈹 Change the language I speak.")

	longDescription = i18n.NewFallbackConfig(
		"plugin.conf.language.long_description",
		"You can list all supported languages by providing no arguments. "+
			"If you want to change my language, give me the language identifier as argument.")

	exampleArgs = []*i18n.Config{
		i18n.EmptyConfig,
		i18n.NewFallbackConfig("plugin.conf.language.example_args.change_language", "en"),
	}
)

// =============================================================================
// Arguments
// =====================================================================================

var argNewLanguageName = i18n.NewFallbackConfig("plugin.conf.language.arg.new_language.name", "New Language")

// =============================================================================
// Responses
// =====================================================================================

var (
	responseList = i18n.NewFallbackConfig(
		"plugin.conf.language.response.list", "You can choose from the following languages: {{.languages}}.")

	responseLanguageChanged = i18n.NewFallbackConfig(
		"plugin.conf.language.response.language_changed",
		"🔧 Successfully changed the language.")
)

type responseListPlaceholders struct {
	Languages string
}

// =============================================================================
// Errors
// =====================================================================================

var errorInvalidLanguage = i18n.NewFallbackConfig(
	"plugin.conf.language.error.invalid_language",
	"`{{.raw}}` isn't one of the languages I speak. You can list all languages by using `{{.language_invoke}}`.")

type errorInvalidLanguagePlaceholders struct {
	Raw            string
	LanguageInvoke string
}
