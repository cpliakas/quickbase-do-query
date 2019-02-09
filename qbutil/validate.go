package qbutil

import (
	"github.com/cpliakas/quickbase-do-query/cliutil"
)

// RequireAppID ensures the table-id option is set.
func RequireAppID(cfg GlobalConfig) string {
	return cliutil.RequireOption(cfg.TableID(), "app-id")
}

// RequireTableID ensures the table-id option is set.
func RequireTableID(cfg GlobalConfig) string {
	return cliutil.RequireOption(cfg.TableID(), "table-id")
}
