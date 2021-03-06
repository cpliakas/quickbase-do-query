package cliutil

import (
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// InitConfig returns a *viper.Viper with an environment variable prefix set
// so that options can be passed via environment variables.
func InitConfig(prefix string) (c *viper.Viper) {
	c = viper.New()
	c.SetEnvPrefix(prefix)
	c.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	c.AutomaticEnv()
	return
}

// Flagger adds flags to a command.
type Flagger struct {
	cmd *cobra.Command
	cfg *viper.Viper
}

// NewFlagger returns a Flagger with
func NewFlagger(cmd *cobra.Command, cfg *viper.Viper) *Flagger {
	return &Flagger{cmd: cmd, cfg: cfg}
}

// Bool adds a local flag that accepts a boolean.
func (f *Flagger) Bool(name, shorthand string, value bool, usage string) {
	f.cmd.Flags().BoolP(name, shorthand, value, usage)
	f.cfg.BindPFlag(name, f.cmd.Flags().Lookup(name))
}

// PersistentBool adds a persistent flag that accepts a boolean.
func (f *Flagger) PersistentBool(name, shorthand string, value bool, usage string) {
	f.cmd.PersistentFlags().BoolP(name, shorthand, value, usage)
	f.cfg.BindPFlag(name, f.cmd.PersistentFlags().Lookup(name))
}

// Int adds a local flag that accepts an integer.
func (f *Flagger) Int(name, shorthand string, value int, usage string) {
	f.cmd.Flags().IntP(name, shorthand, value, usage)
	f.cfg.BindPFlag(name, f.cmd.Flags().Lookup(name))
}

// PersistentInt adds a persistent flag that accepts an integer.
func (f *Flagger) PersistentInt(name, shorthand string, value int, usage string) {
	f.cmd.PersistentFlags().IntP(name, shorthand, value, usage)
	f.cfg.BindPFlag(name, f.cmd.PersistentFlags().Lookup(name))
}

// String adds a local flag that accepts an string.
func (f *Flagger) String(name, shorthand, value, usage string) {
	f.cmd.Flags().StringP(name, shorthand, value, usage)
	f.cfg.BindPFlag(name, f.cmd.Flags().Lookup(name))
}

// PersistentString adds a persistent flag that accepts an string.
func (f *Flagger) PersistentString(name, shorthand, value, usage string) {
	f.cmd.PersistentFlags().StringP(name, shorthand, value, usage)
	f.cfg.BindPFlag(name, f.cmd.PersistentFlags().Lookup(name))
}
