package engine

import (
	"errors"
	"fmt"
	"os"

	"github.com/juliengk/go-cert/pkix"
	"github.com/kassisol/tsa/client"
	"github.com/kassisol/tsa/pkg/adf"
	"github.com/kassisol/twic/pkg/input"
	"github.com/kassisol/twic/pkg/sysutil"
	"github.com/kassisol/twic/pkg/urlutil"
	"github.com/spf13/cobra"
)

func newRemoveCommand() *cobra.Command {
	var (
		tsaToken    string
		tsaUsername string
		tsaPassword string
	)

	cmd := &cobra.Command{
		Use:     "rm",
		Aliases: []string{"remove"},
		Short:   "Remove Docker engine certificate",
		Long:    removeDescription,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runRemove(args, tsaToken, tsaUsername, tsaPassword)
		},
	}

	flags := cmd.Flags()
	flags.StringVarP(&tsaToken, "token", "t", "", "Token")
	flags.StringVarP(&tsaUsername, "username", "u", "", "Username")
	flags.StringVarP(&tsaPassword, "password", "p", "", "Password")

	return cmd
}

func runRemove(args []string, tsaToken, tsaUsername, tsaPassword string) error {
	if !sysutil.IsRoot() {
		return errors.New("you must be root to run engine subcommand")
	}

	if len(args) > 0 {
		return errors.New("this command takes no arguments")
	}

	var username, password string
	if len(tsaToken) == 0 {
		username = tsaUsername
		if len(username) == 0 {
			username = input.ReadPassword("Username")
		}

		password = tsaPassword
		if len(password) == 0 {
			password = input.ReadPassword("Password")
		}
	}

	if len(tsaToken) == 0 {
		if len(username) == 0 {
			return errors.New("empty username is not allowed")
		}
		if len(password) == 0 {
			return errors.New("empty password is not allowed")
		}
	}

	cfg := adf.NewEngine()
	if err := cfg.Init(); err != nil {
		return err
	}

	certificate, err := pkix.NewCertificateFromPEMFile(cfg.TLS.CrtFile)
	if err != nil {
		return err
	}

	crldp := certificate.Crt.CRLDistributionPoints[0]

	url, err := urlutil.Parse(crldp)
	if err != nil {
		return err
	}

	tsaurl := fmt.Sprintf("%s://%s", url.Scheme, url.Host)
	if url.Port != 443 {
		tsaurl = fmt.Sprintf("%s://%s:%d", url.Scheme, url.Host, url.Port)
	}

	clt, err := client.New(tsaurl)
	if err != nil {
		return err
	}

	if err := clt.GetDirectory(); err != nil {
		return err
	}

	token := tsaToken
	if len(token) == 0 {
		token, err = clt.GetToken(username, password, 0)
		if err != nil {
			return err
		}
	}

	if err := clt.RevokeCertificate(token, int(certificate.Crt.SerialNumber.Int64())); err != nil {
		return err
	}

	if err := os.Remove(cfg.TLS.CaFile); err != nil {
		return err
	}

	if err := os.Remove(cfg.TLS.KeyFile); err != nil {
		return err
	}

	return os.Remove(cfg.TLS.CrtFile)
}

var removeDescription = `
Remove Docker engine certificate

`
