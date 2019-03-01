package qbutil

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/cpliakas/quickbase-do-query/cliutil"
	"github.com/cpliakas/quickbase-do-query/qb"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"github.com/spf13/viper"
)

// GlobalConfig contains configuration common to all commands.
type GlobalConfig struct {
	viper *viper.Viper

	// RequireTableID flags that a table ID is required by the command.
	RequireTableID bool

	// RequireAppID flags that an app ID is required by the command.
	RequireAppID bool

	// RequireCreds flags that credentials are required by the command.
	RequireCreds bool
}

// NewGlobalConfig returns a GlobalConfig.
func NewGlobalConfig(cmd *cobra.Command, cfg *viper.Viper) GlobalConfig {
	flags := cliutil.NewFlagger(cmd, cfg)

	flags.PersistentString("app-id", "I", "", "application's dbid")
	flags.PersistentString("app-token", "A", "", "app token used with ticket to to authenticate API requests")
	flags.PersistentBool("batch", "B", false, "render output in batch mode, useful for chaining commands together")
	flags.PersistentString("config-file", "C", qb.DefaultConfigFile, "path to the config file")
	flags.PersistentString("filter", "F", "", "JMESPath filter")
	flags.PersistentBool("raw", "X", false, "return the raw output from the API call")
	flags.PersistentString("realm-host", "R", "", "realm host, e.g., 'https://MYREALM.quickbase.com'")
	flags.PersistentString("table-id", "t", "", "table's dbid")
	flags.PersistentString("ticket", "T", "", "ticket used to authenticate API requests")
	flags.PersistentString("ticket-file", "K", qb.DefaultTicketFile, "path to the file containing a cached ticket")
	flags.PersistentString("user-token", "U", "", "user token used to authenticate API requests")

	return GlobalConfig{viper: cfg}
}

// Set implements qb.Config.Set.
func (c GlobalConfig) Set(key string, value interface{}) {
	c.viper.Set(key, value)
}

// AppID implements qb.Config.AppID.
func (c GlobalConfig) AppID() string { return c.viper.GetString("app-id") }

// AppToken implements qb.Config.AppToken.
func (c GlobalConfig) AppToken() string { return c.viper.GetString("app-token") }

// Batch returns whether to render output in batch mode.
func (c GlobalConfig) Batch() bool { return c.viper.GetBool("batch") }

// ConfigFile implements qb.Config.ConfigFile.
func (c GlobalConfig) ConfigFile() string { return c.viper.GetString("config-file") }

// Filter returns the JMESPath filter.
func (c GlobalConfig) Filter() string { return c.viper.GetString("filter") }

// Raw flags whether to return the raw output from the API as opposed to JSON.
func (c GlobalConfig) Raw() bool { return c.viper.GetBool("raw") }

// RealmHost implements qb.Config.RealmHost.
func (c GlobalConfig) RealmHost() string { return c.viper.GetString("realm-host") }

// TableID returns the configured table's dbid.
func (c GlobalConfig) TableID() string { return c.viper.GetString("table-id") }

// Ticket implements qb.Config.Ticket.
func (c GlobalConfig) Ticket() string { return c.viper.GetString("ticket") }

// TicketFile implements qb.Config.TicketFile.
func (c GlobalConfig) TicketFile() string { return c.viper.GetString("ticket-file") }

// UserToken implements qb.Config.UserToken.
func (c GlobalConfig) UserToken() string { return c.viper.GetString("user-token") }

// InitConfig wraps qb.InitConfig.
func (c *GlobalConfig) InitConfig() error {
	if err := qb.InitConfig(c.viper); err != nil {
		return fmt.Errorf("error reading configuration: %s", err)
	}
	return nil
}

// Validate validates the global configuration options.
func (c *GlobalConfig) Validate() error {
	if err := c.InitConfig(); err != nil {
		return err
	}

	// The realm-host option is always required.
	if err := validation.Validate(c.RealmHost(),
		validation.Required,
		is.URL,
	); err != nil {
		return fmt.Errorf("realm-host option invalid: %s", err)
	}

	// Validate the app-id option.
	if c.RequireTableID {
		if err := validation.Validate(c.AppID(),
			validation.Required,
			validation.Length(9, 9),
		); err != nil {
			return fmt.Errorf("app-id option invalid: %s", err)
		}
	}

	// Validate the table-id option.
	if c.RequireTableID {
		if err := validation.Validate(c.TableID(),
			validation.Required,
			validation.Length(9, 9),
		); err != nil {
			return fmt.Errorf("table-id option invalid: %s", err)
		}
	}

	return nil
}
