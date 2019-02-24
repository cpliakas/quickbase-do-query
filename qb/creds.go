package qb

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// TicketHours is the default expiry for a ticket.
const TicketHours = 12

// Credentials contains the parameters that are used to authenticate
// requests made to the Quick Base API.
// See https://quickbase.com/api-guide/authentication_and_secure_access.html.
type Credentials struct {
	AppToken  string `xml:"apptoken,omitempty"`
	Ticket    string `xml:"ticket,omitempty"`
	UserToken string `xml:"usertoken,omitempty"`
}

// NewCredentials returns a Credentials with data populated dependant on
// configuration.
func NewCredentials(cfg Config) Credentials {
	if cfg.UserToken() != "" {
		return Credentials{UserToken: cfg.UserToken()}
	}

	if cfg.Ticket() != "" {
		if cfg.AppToken() == "" {
			return Credentials{Ticket: cfg.Ticket()}
		}
		return Credentials{Ticket: cfg.Ticket(), AppToken: cfg.AppToken()}
	}

	return Credentials{}
}

// NewTicket calls the API_Authenticate endpoint to create a ticket that can
// be used to authenticate subsequent API requests. If a ticket is returned,
// it is cached in the ticket file.
func (c Client) NewTicket(cachefile, username, password string, hours int) (output AuthenticateOutput, err error) {
	input := &AuthenticateInput{
		Username: username,
		Password: password,
		Hours:    hours,
	}

	output, err = c.Authenticate(input)
	if err != nil {
		return
	}

	err = CacheTicket(cachefile, output)
	if err != nil {
		err = TicketFileError{cachefile, err}
	}

	return
}

// CacheTicket writes the ticket to file so that it can be used in
// subsequent API requests.
func CacheTicket(file string, output AuthenticateOutput) error {
	dir := filepath.Dir(file)

	// Ensure the cache directory exists.
	stat, err := os.Stat(dir)
	if os.IsNotExist(err) {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return err
		}
	} else if err != nil {
		return err
	} else if !stat.IsDir() {
		return err
	}

	// Write the ticket to the ticket file.
	b := []byte(fmt.Sprintf("%s\n", output.Ticket))
	return ioutil.WriteFile(file, b, 0600)
}

// ReadCachedTicket reads the data in the ticket file and returns the ticket
// if it exists.
func ReadCachedTicket(file string) (string, error) {

	// Ensure the file exists.
	stat, err := os.Stat(file)
	if os.IsNotExist(err) {
		return "", nil
	} else if err != nil {
		return "", TicketFileError{file, err}
	} else if stat.IsDir() {
		return "", errors.New("cache file is a directory, expected a regular file")
	}

	// Read the contents of the file.
	b, err := ioutil.ReadFile(file)
	if err != nil {
		return "", err
	}

	// Return the cached ticket.
	// TODO validate the ticket.
	return strings.TrimRight(string(b), "\n"), nil
}

// TicketFileError implents the error interface and records an error reading
// from and writing to the ticket file.
type TicketFileError struct {
	File string
	Err  error
}

// Error satisfies the error interface and simply returns the error string
// of the underlying error.
func (e TicketFileError) Error() string { return e.Err.Error() }

// IsTicketFileErr returns true if the error is a TicketFileError.
func IsTicketFileErr(err error) bool {
	_, ok := err.(TicketFileError)
	return ok
}
