package cert

import (
	"errors"
	"fmt"
	"os"

	"github.com/juliengk/go-cert/helpers"
	"github.com/juliengk/go-cert/pkix"
	"github.com/kassisol/tsa/client"
	"github.com/kassisol/tsa/pkg/adf"
	"github.com/kassisol/twic/pkg/input"
	"github.com/kassisol/twic/storage"
	"github.com/spf13/cobra"
)

func newRenewCommand() *cobra.Command {
	var (
		tsaToken    string
		tsaPassword string
	)

	cmd := &cobra.Command{
		Use:   "renew [name]",
		Short: "Renew Docker client certificate",
		Long:  renewDescription,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runRenew(args, tsaToken, tsaPassword)
		},
	}

	flags := cmd.Flags()
	flags.StringVarP(&tsaToken, "token", "t", "", "Token")
	flags.StringVarP(&tsaPassword, "password", "p", "", "Password")

	return cmd
}

func runRenew(args []string, tsaToken, tsaPassword string) error {
	if len(args) != 1 {
		return errors.New("this command requires exactly one argument")
	}

	name := args[0]

	cfg := adf.NewClient()
	if err := cfg.Init(); err != nil {
		return err
	}

	cfg.SetName(name)

	s, err := storage.NewDriver("sqlite", cfg.App.Dir.Root)
	if err != nil {
		return err
	}
	defer s.End()

	crt, err := s.GetCert(name)
	if err != nil {
		return fmt.Errorf("name %q does not exist", name)
	}

	clt, err := client.New(crt.TSAURL)
	if err != nil {
		return err
	}

	if err := clt.GetDirectory(); err != nil {
		return err
	}

	oldcert, err := pkix.NewCertificateFromPEMFile(cfg.TLS.CrtFile)
	if err != nil {
		return err
	}

	key, err := pkix.NewKey(4096)
	if err != nil {
		return err
	}

	keyBytes, err := key.ToPEM()
	if err != nil {
		return err
	}

	csr, err := helpers.CreateCSR(oldcert.Crt.Subject.Country[0], oldcert.Crt.Subject.Province[0], oldcert.Crt.Subject.Locality[0], oldcert.Crt.Subject.Organization[0], oldcert.Crt.Subject.OrganizationalUnit[0], oldcert.Crt.Subject.CommonName, "", []string{}, key)
	if err != nil {
		return err
	}

	token := tsaToken
	if len(token) == 0 {
		password := tsaPassword
		if len(password) == 0 {
			password = input.ReadPassword("Password")
		}
		token, err = clt.GetToken(crt.CN, password, 0)
		if err != nil {
			return err
		}
	}

	if err := clt.RevokeCertificate(token, int(oldcert.Crt.SerialNumber.Int64())); err != nil {
		return err
	}

	newcert, err := clt.GetCertificate(token, "client", csr.Bytes, 12)
	if err != nil {
		return err
	}

	if err := os.Remove(cfg.TLS.CrtFile); err != nil {
		return err
	}

	if err := os.Remove(cfg.TLS.KeyFile); err != nil && !os.IsNotExist(err) {
		return err
	}

	if err := pkix.ToPEMFile(cfg.TLS.CrtFile, []byte(newcert), 0444); err != nil {
		return err
	}

	return pkix.ToPEMFile(cfg.TLS.KeyFile, keyBytes, 0400)
}

var renewDescription = `
Renew Docker client certificate

`
