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

var argNewPrefixName = i18n.NewFallbackConfig("plugin.conf.prefix.arg.plugin.name", "New Prefix")

// =============================================================================
// Flags
// =====================================================================================

var flagRemoveDescription = i18n.NewFallbackConfig(
	"plugin.conf.prefix.flag.remove.description",
	"If you set this flag, the current prefix will be removed without replacement. "+
		"You can still use my commands by mentioning me.")

// =============================================================================
// Responses
// =====================================================================================

var (
	responseCurrentPrefix = i18n.NewFallbackConfig(
		"plugin.conf.prefix.response.current_prefix",
		"The prefix on this server is ``{{.prefix}}``.")

	responseNoPrefix = i18n.NewFallbackConfig(
		"plugin.conf.prefix.response.no_prefix",
		"There is no custom prefix on this server.")

	responsePrefixChanged = i18n.NewFallbackConfig(
		"plugin.conf.prefix.response.prefix_changed",
		"🔧 Successfully changed the prefix to ``{{.new_prefix}}``.")
)

type (
	responseCurrentPrefixPlaceholders struct {
		Prefix string
	}

	responsePrefixChangedPlaceholders struct {
		NewPrefix string
	}
)
