package cliutil

import (
	"errors"
	"strconv"
)

// RequireArg ensures that an argument was passed.
func RequireArg(args []string, pos int, name string) string {
	if len(args) < pos+1 {
		HandleError(errors.New(name), "missing required argument")
	}
	return args[pos]
}

// RequireArgInt ensures that an argument was passed and is an integer.
func RequireArgInt(args []string, pos int, name string) int {
	s := RequireArg(args, pos, name)
	i, err := strconv.Atoi(s)
	if err != nil {
		HandleError(errors.New(name), "argument must be an integer")
	}
	return i
}

// RequireOption ensures that an option was passed.
func RequireOption(val, name string) string {
	if val == "" {
		HandleError(errors.New(name), "missing required option")
	}
	return val
}
