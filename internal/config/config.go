// Package config provides the configuration of levin.
package config

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"
	_ "time/tzdata" // used for parsing locations

	"github.com/diamondburned/arikawa/v2/discord"
	"github.com/diamondburned/arikawa/v2/gateway"
	"github.com/iancoleman/strcase"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"golang.org/x/text/language"
)

// C is the global config of levin.
var C config

func log() *zap.SugaredLogger { return zap.S().Named("config") }

// config is the type holding all configurational data.
type config struct { //nolint:maligned
	Token  string `mapstructure:"bot_token"`
	Owners []discord.UserID

	DefaultPrefix   string         `mapstructure:"default_prefix"`
	DefaultLanguage language.Tag   `mapstructure:"-"`
	DefaultTimeZone *time.Location `mapstructure:"-"`

	Status       gateway.Status
	ActivityType discord.ActivityType `mapstructure:"-"`
	ActivityName string               `mapstructure:"-"`
	ActivityURL  discord.URL

	EditAge  time.Duration `mapstructure:"edit_age"`
	AllowBot bool          `mapstructure:"allow_bot"`

	Sentry struct {
		DSN         string
		Environment string
		SampleRate  float64 `mapstructure:"sample_rate"`
	}

	Mongo struct {
		URI          string
		DatabaseName string `mapstructure:"database_name"`
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

	v.RegisterAlias("mongo.db_name", "mongo.database_name")
	_ = v.BindEnv("mongo.database_name", "LEVIN_MONGO_DB_NAME")

	v.RegisterAlias("default_lang", "default_language")
	_ = v.BindEnv("default_language", "LEVIN_DEFAULT_LANG")

	v.RegisterAlias("default_tz", "default_time_zone")
	_ = v.BindEnv("default_time_zone", "LEVIN_DEFAULT_TZ")

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
			if err := bindEnvs(v, f.Type, name+"."); err != nil {
				return err
			}
		} else {
			envName := strings.ReplaceAll(strings.ToUpper(name), ".", "_")

			err := v.BindEnv(name, fmt.Sprintf("LEVIN_%s", envName))
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

	C.DefaultLanguage = parseLanguage(v.GetString("default_language"))
	C.DefaultTimeZone = parseTimeZone(v.GetString("default_time_zone"))
	C.EditAge = time.Duration(v.GetInt("edit_age")) * time.Second
	C.ActivityType, C.ActivityName = parseActivity(v.GetString("activity"))
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

func parseLanguage(lang string) language.Tag {
	t, err := language.Parse(lang)
	if err != nil {
		return language.English
	}

	return t
}

func parseTimeZone(tzstr string) *time.Location {
	tz, err := time.LoadLocation(tzstr)
	if err != nil {
		return time.UTC
	}

	return tz
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
