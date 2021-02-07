// Package config provides the configuration of levin.
package config

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/diamondburned/arikawa/v2/discord"
	"github.com/diamondburned/arikawa/v2/gateway"
	"github.com/iancoleman/strcase"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

// C is the global config of levin.
var C config

func log() *zap.SugaredLogger { return zap.S().Named("config") }

// config is the type holding all configurational data.
type config struct {
	Token           string   `mapstructure:"bot_token"`
	DefaultPrefixes []string `mapstructure:"default_prefixes"`
	Owners          []discord.UserID

	Status       gateway.Status
	ActivityType discord.ActivityType `mapstructure:"-"`
	ActivtyName  string               `mapstructure:"-"`

	EditAge  time.Duration `mapstructure:"edit_age"`
	AllowBot bool          `mapstructure:"allow_bot"`

	Sentry struct {
		DSN         string
		SampleRate  float64 `mapstructure:"sample_rate"`
		Environment string
	}

	ServerName string `mapstructure:"server_name"`
}

// Zero sets all config fields to their zero values.
func Zero() { C = config{} }

// Load loads the config.
// If configPath is not empty, the config at that path will be loaded, instead
// of searching in the current directory, ./config and $CONFIG_DIR/
func Load(configPath string) error {
	v := viper.New()

	v.SetEnvPrefix("levin")
	v.AutomaticEnv()

	if err := bindEnvs(v, reflect.TypeOf(C), ""); err != nil {
		return err
	}

	if len(configPath) > 0 {
		v.SetConfigFile(configPath)
	} else {
		v.AddConfigPath(".")
		v.AddConfigPath("config")
		v.AddConfigPath("$CONFIG_DIR/")
		v.SetConfigName("levin")
	}

	loadDefaults(v)

	err := v.ReadInConfig()
	if err != nil && (!errors.As(err, new(viper.ConfigFileNotFoundError)) || len(configPath) > 0) {
		return err
	}

	err = unmarshal(v)
	if err == nil {
		log().With("config", C).
			Debug("read config")
	}

	return err
}

// bindEnvs binds all config fields to environment variables.
// This is necessary, because Viper.AutomaticEnv() does not work for
// unmarshalling into structs
func bindEnvs(v *viper.Viper, val reflect.Type, base string) error {
	for i := 0; i < val.NumField(); i++ {
		f := val.Field(i)
		if f.Anonymous { // viper can't set it, so we don't care
			continue
		}

		name := f.Tag.Get("mapstructure")
		if name == "-" { // skip
			continue
		}

		if len(name) == 0 {
			name = strcase.ToSnake(f.Name)
		}

		name = base + name

		if f.Type.Kind() == reflect.Struct {
			if err := bindEnvs(v, f.Type, name); err != nil {
				return err
			}
		} else {
			err := v.BindEnv(name, fmt.Sprintf("LEVIN_%s", strings.ToUpper(name)))
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func loadDefaults(v *viper.Viper) {
	v.SetDefault("allow_bot", false)
	v.SetDefault("edit_age", 15 /* seconds */)
}

func unmarshal(v *viper.Viper) error {
	if err := v.Unmarshal(&C); err != nil {
		return err
	}

	C.EditAge = time.Duration(v.GetInt("edit_age")) * time.Second
	C.ActivityType, C.ActivtyName = parseActivity(v.GetString("activity"))
	C.Status = validateStatus(gateway.Status(v.GetString("status")))

	return nil
}

func validateStatus(status gateway.Status) gateway.Status {
	switch status {
	case gateway.UnknownStatus:
	case gateway.OnlineStatus:
	case gateway.DoNotDisturbStatus:
	case gateway.IdleStatus:
	case gateway.OfflineStatus:
		status = gateway.InvisibleStatus
	case gateway.InvisibleStatus:
	default:
		status = gateway.OnlineStatus
	}

	return status
}

var activityTypes = map[discord.ActivityType]string{
	discord.GameActivity:      "Playing",
	discord.StreamingActivity: "Streaming",
	discord.ListeningActivity: "Listening to",
	discord.WatchingActivity:  "Watching",
}

func parseActivity(activity string) (t discord.ActivityType, name string) {
	for t, tstr := range activityTypes {
		if strings.HasPrefix(activity, tstr) {
			if len(activity) > len(tstr)+1 {
				return t, activity[len(tstr)+1:]
			}

			return 0, ""
		}
	}

	return 0, ""
}
