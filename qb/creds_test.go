package qb

import (
	"errors"
	"io/ioutil"
	"net/http"
	"os"
	"testing"
)

const mockAppToken = "cvukcpupfz3a3ed6qyctybymhwy9"
const mockTicket = "8_bm58qftf6_b1b6se_uyp_a_dnqpixytb64vqc95hvs2b289tnt4xfi5qcub9gtbb6qnpaqqzzurx"
const mockUserID = "12345678.abcd"
const mockUserToken = "b3b7se_uvp_b36r2kdd6fu3q3ftm8mkcbfk1uzc"

func TempDir(t *testing.T) string {
	dir, err := ioutil.TempDir(os.TempDir(), "quickbase-sdk-")
	if err != nil {
		t.Fatalf("error creating temp dir: %s", err)
	}
	return dir
}

func TempFile(t *testing.T) string {
	file, err := ioutil.TempFile(os.TempDir(), "quickbase-sdk-")
	if err != nil {
		t.Fatalf("error creating temp file: %s", err)
	}
	return file.Name()
}

func TestIsTicketFileErr(t *testing.T) {
	err := errors.New("test error")
	tfe := TicketFileError{"/dev/null", err}

	if !IsTicketFileErr(tfe) {
		t.Error("expected error to be a TicketFileError")
	}

	if IsTicketFileErr(err) {
		t.Error("expected error to not be a TicketFileError")
	}
}

func TestTicketFileErr(t *testing.T) {
	err := errors.New("test error")
	tfe := TicketFileError{"/dev/null", err}

	if tfe.Error() != err.Error() {
		t.Errorf("expected '%s', got '%s'", err, tfe)
	}
}

func TestNewCredentialsUserToken(t *testing.T) {
	cfg := newTestConfig()
	cfg.viper.Set("user-token", mockUserToken)
	cfg.viper.Set("ticket", mockTicket)

	c := NewCredentials(cfg)
	if c.UserToken != mockUserToken {
		t.Errorf("expected '%s', got '%s'", mockUserToken, c.UserToken)
	}
	if c.Ticket != "" {
		t.Error("expected user token to take precedence over ticket")
	}
}

func TestNewCredentialsTicket(t *testing.T) {
	cfg := newTestConfig()
	cfg.viper.Set("ticket", mockTicket)

	c := NewCredentials(cfg)
	if c.Ticket != mockTicket {
		t.Errorf("expected '%s', got '%s'", mockTicket, c.Ticket)
	}
}

func TestNewCredentialsTicketWithAppToken(t *testing.T) {
	cfg := newTestConfig()
	cfg.viper.Set("ticket", mockTicket)
	cfg.viper.Set("app-token", mockAppToken)

	c := NewCredentials(cfg)
	if c.Ticket != mockTicket {
		t.Errorf("expected ticket '%s', got '%s'", mockTicket, c.Ticket)
	}
	if c.AppToken != mockAppToken {
		t.Errorf("expected app token '%s', got '%s'", mockAppToken, c.AppToken)
	}
}

func TestNewCredentialsWithoutCreds(t *testing.T) {
	cfg := newTestConfig()

	c := NewCredentials(cfg)
	if c.AppToken != "" {
		t.Errorf("expected empty app token, got '%s'", c.AppToken)
	}
	if c.Ticket != "" {
		t.Errorf("expected empty ticket got '%s'", c.Ticket)
	}
	if c.UserToken != "" {
		t.Errorf("expected empty user token, got '%s'", c.UserToken)
	}
}

func authenticateSuccessHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`<?xml version="1.0" ?>
		<qdbapi>
			<action>API_Authenticate</action>
			<errcode>0</errcode>
			<errtext>No error</errtext>
			<ticket>` + mockTicket + `</ticket>
			<userid>` + mockUserID + `</userid>
		</qdbapi>`))
}

func TestNewTicket(t *testing.T) {

	server, client := NewServerClientPair(authenticateSuccessHandler)
	defer server.Close()

	file := TempFile(t)
	defer os.Remove(file)

	out, err := client.NewTicket(file, "username", "password", 4)
	if err != nil {
		t.Fatalf("error requesting new ticket: %s", err)
	}

	if out.Ticket != mockTicket {
		t.Errorf("expected returned ticket '%s', got '%s'", mockTicket, out.Ticket)
	}
	if out.UserID != mockUserID {
		t.Errorf("expected returned user id '%s', got '%s'", mockUserID, out.UserID)
	}

	ticket, err := ReadCachedTicket(file)
	if err != nil {
		t.Fatalf("error reading ticket file: %s", err)
	}

	if ticket != out.Ticket {
		t.Errorf("expected cached ticket '%s', got '%s'", out.Ticket, ticket)
	}
}

func authenticateBadCredsHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`<?xml version="1.0" ?>
		<qdbapi>
			<action>API_Authenticate</action>
			<errcode>20</errcode>
			<errtext>Unknown username/password</errtext>
			<errdetail>Sorry! You entered the wrong E-Mail or Screen Name or Password. Try again.</errdetail>
		</qdbapi>`))
}

func TestNewTicketBadCreds(t *testing.T) {

	server, client := NewServerClientPair(authenticateBadCredsHandler)
	defer server.Close()

	file := TempFile(t)
	defer os.Remove(file)

	out, err := client.NewTicket(file, "username", "password", 4)
	if err == nil {
		t.Fatal("expected error requesting new ticket")
	}
	if out.ErrorCode != 20 {
		t.Fatalf("expected error code '20', got '%v'", out.ErrorCode)
	}
}
