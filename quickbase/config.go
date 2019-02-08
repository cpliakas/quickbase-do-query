package quickbase

import (
	"fmt"
	"os"
	"strings"

	"github.com/kubernetes/client-go/util/homedir"
	"github.com/spf13/viper"
)

// Config is the interface that contains runtime configuration.
type Config interface {

	// AppID returns the dbid of the application used by default.
	AppID() string

	// AppToken returns the app token used to authenticate API requests.
	// See https://quickbase.com/api-guide/authentication_and_secure_access.html.
	AppToken() string

	// ConfigFile retuns the path to the TOML file that contains config. The
	// default path is "$HOME/.config/quickbase/config", where $HOME is
	// replaced with the user's home directory dependant on platform.
	ConfigFile() string

	// The realm host, e.g., https://MYREALM.quickbase.com/. Replace MYREALM
	// accordingly.
	RealmHost() string

	// Ticket returns the ticket used to authenticate API requests.
	// See https://quickbase.com/api-guide/authentication_and_secure_access.html
	Ticket() string

	// TicketFile returns the path to file continaing a cached ticket. The
	// default path is "$HOME/.config/quickbase/ticket", where $HOME is
	// replaced with the user's home directory dependant on platform.
	TicketFile() string

	// UserToken returns the user token used to authenticate API requests.
	// See https://quickbase.com/api-guide/create_user_tokens.html.
	UserToken() string
}

// NewConfig returns a StandardConfig that reads configuration from
// explicitly set options, environment variables, and a TOML configuration
// file in that order of preference.
func NewConfig() StandardConfig {
	v := viper.New()

	// Read from QUICK_BASE_ environment Variables.
	v.SetEnvPrefix("QUICKBASE")
	v.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	v.AutomaticEnv()

	// Bind environment variables.
	v.BindEnv("app-id")
	v.BindEnv("app-token")
	v.BindEnv("config-file")
	v.BindEnv("realm-host")
	v.BindEnv("ticket")
	v.BindEnv("ticket-file")
	v.BindEnv("user-token")

	// Set defaults.
	v.SetDefault("config-file", "$HOME/.config/quickbase/config")
	v.SetDefault("host", "quickbase.com")
	v.SetDefault("ticket-file", "$HOME/.config/quickbase/ticket")
	v.SetDefault("scheme", "https")

	// Read in the config file.
	err := InitConfig(v)
	if err != nil {
		panic(fmt.Sprintf("error reading config file: %s", err))
	}

	return StandardConfig{viper: v}
}

// InitConfig reads config from options, environment variables, and config
// files in that order of preference. This is separated out into a different
// function so that it can be used by methods that construct their own
// instance of *viper.Viper.
func InitConfig(v *viper.Viper) error {

	// Read ticket from ticket file if one isn't already set.
	v.Set("ticket-file", ReplaceTokens(v, "ticket-file"))
	if v.GetString("ticket") == "" {
		ticket, err := ReadCachedTicket(v.GetString("ticket-file"))
		if err == nil && ticket != "" {
			v.Set("ticket", ticket)
		}
	}

	// Read configuration from the configuration file if it exists.
	configFile := ReplaceTokens(v, "config-file")
	_, err := os.Stat(configFile)
	if !os.IsNotExist(err) {
		v.SetConfigFile(configFile)
		v.SetConfigType("toml")
		return v.ReadInConfig()
	}

	return nil
}

// ReplaceTokens replaces tokens in the values set as configuration options.
// For example, $HOME is replaced with the user's home directory dependant
// on platform, see https://github.com/kubernetes/client-go/blob/master/util/homedir/homedir.go.
func ReplaceTokens(v *viper.Viper, key string) string {
	return strings.Replace(v.GetString(key), "$HOME", homedir.HomeDir(), -1)
}

// StandardConfig wraps viper.Viper "Get*" methods to return configuration.
type StandardConfig struct {
	viper *viper.Viper
}

// AppID implements Config.AppID().
func (c StandardConfig) AppID() string { return c.viper.GetString("app-id") }

// AppToken implements Config.AppToken().
func (c StandardConfig) AppToken() string { return c.viper.GetString("app-token") }

// ConfigFile implements Config.ConfigFile().
func (c StandardConfig) ConfigFile() string { return c.viper.GetString("config-file") }

// RealmHost implements Config.RealmHost().
func (c StandardConfig) RealmHost() string { return c.viper.GetString("realm-host") }

// Ticket implements Config.Ticket().
func (c StandardConfig) Ticket() string { return c.viper.GetString("ticket") }

// TicketFile implements Config.TicketFile().
func (c StandardConfig) TicketFile() string { return c.viper.GetString("ticket-file") }

// UserToken implements Config.UserToken().
func (c StandardConfig) UserToken() string { return c.viper.GetString("user-token") }
