package botmaster

import "github.com/mavolin/adam/pkg/i18n"

var botMasterError = i18n.NewFallbackConfig(
	"bot_master.error.not_master",
	"You need to be a bot master administrator to use this command")
