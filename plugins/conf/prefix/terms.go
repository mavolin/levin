package prefix

import "github.com/mavolin/adam/pkg/i18n"

// =============================================================================
// Meta
// =====================================================================================

var (
	shortDescription = i18n.NewFallbackConfig(
		"plugin.conf.prefix.short_description",
		"Show and change my prefix.")

	longDescription = i18n.NewFallbackConfig(
		"plugin.conf.prefix.long_description",
		"Display the prefix used on this server, or change it by providing a new one.")

	exampleArgs = []*i18n.Config{
		i18n.EmptyConfig,
		i18n.NewFallbackConfig("plugin.conf.prefix.example_args.new", "!"),
	}
)

// =============================================================================
// Arguments
// =====================================================================================

var (
	argNewPrefixName     = i18n.NewFallbackConfig("plugin.conf.prefix.arg.plugin.name", "New Prefix")
	argPluginDescription = i18n.NewFallbackConfig(
		"plugin.conf.prefix.arg.plugin.description",
		"The new prefix you want me to use. It may not contain whitespace.")
)
