package main

import (
	"errors"
	"time"
)

type (
	cliConfig struct {
		ClientID       string        `flag:"client-id" default:"" description:"Client-ID of Mercedes Developers Console App"`
		ClientSecret   string        `flag:"client-secret" default:"" description:"Client-Secret of Mercedes Developers Console App"`
		CredentialFile string        `flag:"credential-file" default:"credentials.json" description:"Where to store tokens when using client-id from CLI parameters"`
		FetchInterval  time.Duration `flag:"fetch-interval" default:"15m" description:"How often to ask the Mercedes API for updates"`
		InfluxExport   string        `flag:"influx-export" default:"" description:"Set to url (http[s]://user:pass@host[:port]/database) to enable Influx exporter"`
		Listen         string        `flag:"listen" default:":3000" description:"Port/IP to listen on"`
		LogLevel       string        `flag:"log-level" default:"info" description:"Log level (debug, info, warn, error, fatal)"`
		RedirectURL    string        `flag:"redirect-url" default:"http://127.0.0.1:3000/store-token" description:"Redirect URL registered in Mercedes Developers Console"`
		VaultKey       string        `flag:"vault-key" default:"" description:"Use credentials from and update in Vault"`
		VehicleID      []string      `flag:"vehicle-id" default:"" description:"Vehicle identification number (e.g. WDB111111ZZZ22222)"`
		VersionAndExit bool          `flag:"version" default:"false" description:"Prints current version and exits"`
	}
)

func (c cliConfig) Validate() error {
	switch {
	case len(c.VehicleID) == 0:
		return errors.New("at least one vehicle-id is required")

	case c.VaultKey == "" && c.ClientID == "":
		return errors.New("either vault-key or client-id/secret is required")

	case c.ClientID != "" && c.ClientSecret == "":
		return errors.New("client-id is set and client-secret is not")

	case c.ClientID != "" && c.VaultKey != "":
		return errors.New("client-id and vault-key are configured, use only one of them")

	default:
		// No errors
		return nil
	}
}
