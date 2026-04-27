package cert

import (
	"errors"
	"fmt"
	"os"

	"github.com/juliengk/go-cert/pkix"
	"github.com/kassisol/tsa/client"
	"github.com/kassisol/tsa/pkg/adf"
	"github.com/kassisol/twic/pkg/input"
	"github.com/kassisol/twic/pkg/sysutil"
	"github.com/kassisol/twic/storage"
	"github.com/spf13/cobra"
)

func newRemoveCommand() *cobra.Command {
	var (
		tsaToken    string
		tsaPassword string
	)

	cmd := &cobra.Command{
		Use:     "rm [name]",
		Aliases: []string{"remove"},
		Short:   "Remove Docker client certificate",
		Long:    removeDescription,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runRemove(args, tsaToken, tsaPassword)
		},
	}

	flags := cmd.Flags()
	flags.StringVarP(&tsaToken, "token", "t", "", "Token")
	flags.StringVarP(&tsaPassword, "password", "p", "", "Password")

	return cmd
}

func runRemove(args []string, tsaToken, tsaPassword string) error {
	if sysutil.IsRoot() {
		return errors.New("you must not be root to remove a client certificate")
	}

	if len(args) != 1 {
		return errors.New("this command requires exactly one argument")
	}

	cfg := adf.NewClient()
	if err := cfg.Init(); err != nil {
		return err
	}

	cfg.SetName(args[0])

	s, err := storage.NewDriver("sqlite", cfg.App.Dir.Root)
	if err != nil {
		return err
	}
	defer s.End()

	var password string
	if len(tsaToken) == 0 {
		password = tsaPassword
		if len(password) == 0 {
			password = input.ReadPassword("Password")
		}
	}

	crt, err := s.GetCert(args[0])
	if err != nil {
		return fmt.Errorf("name %q does not exist", args[0])
	}

	if len(tsaToken) == 0 && len(password) == 0 {
		return errors.New("empty password is not allowed")
	}

	clt, err := client.New(crt.TSAURL)
	if err != nil {
		return err
	}

	if err := clt.GetDirectory(); err != nil {
		return err
	}

	token := tsaToken
	if len(token) == 0 {
		token, err = clt.GetToken(crt.CN, password, 0)
		if err != nil {
			return err
		}
	}

	certificate, err := pkix.NewCertificateFromPEMFile(cfg.TLS.CrtFile)
	if err != nil {
		return err
	}

	if err := clt.RevokeCertificate(token, int(certificate.Crt.SerialNumber.Int64())); err != nil {
		return err
	}

	if err := s.RemoveCert(args[0]); err != nil {
		return err
	}

	return os.RemoveAll(cfg.Profile.CertDir)
}

var removeDescription = `
Remove Docker client certificate

`
