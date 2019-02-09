package qbutil

import (
	"fmt"

	"github.com/cpliakas/quickbase-do-query/quickbase"
	"github.com/spf13/viper"
)

// GlobalConfig contains configuration common to all commands.
type GlobalConfig struct {
	viper *viper.Viper
}

// NewGlobalConfig returns a GlobalConfig.
func NewGlobalConfig(v *viper.Viper) GlobalConfig {
	return GlobalConfig{viper: v}
}

// InitConfig wraps quickbase.InitConfig().
// TODO Remove the panic?
func (c *GlobalConfig) InitConfig() {
	err := quickbase.InitConfig(c.viper)
	if err != nil {
		panic(fmt.Errorf("error reading config file: %s", err))
	}
}

// AppID implements quickbase.Config.AppID().
func (c GlobalConfig) AppID() string { return c.viper.GetString("app-id") }

// AppToken implements quickbase.Config.AppToken().
func (c GlobalConfig) AppToken() string { return c.viper.GetString("app-token") }

// Batch returns whether to render output in batch mode.
func (c GlobalConfig) Batch() bool { return c.viper.GetBool("batch") }

// ConfigFile implements quickbase.Config.ConfigFile().
func (c GlobalConfig) ConfigFile() string { return c.viper.GetString("config-file") }

// Filter returns the JMESPath filter.
func (c GlobalConfig) Filter() string { return c.viper.GetString("filter") }

// Raw flags whether to return the raw output from the API as opposed to JSON.
func (c GlobalConfig) Raw() bool { return c.viper.GetBool("raw") }

// RealmHost implements quickbase.Config.RealmHost().
func (c GlobalConfig) RealmHost() string { return c.viper.GetString("realm-host") }

// TableID returns the configured table's dbid.
func (c GlobalConfig) TableID() string { return c.viper.GetString("table-id") }

// Ticket implements quickbase.Config.Ticket().
func (c GlobalConfig) Ticket() string { return c.viper.GetString("ticket") }

// TicketFile implements quickbase.Config.TicketFile().
func (c GlobalConfig) TicketFile() string { return c.viper.GetString("ticket-file") }

// UserToken implements quickbase.Config.UserToken().
func (c GlobalConfig) UserToken() string { return c.viper.GetString("user-token") }
