package cert

import (
	"errors"
	"fmt"
	"os"
	"strings"

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
		byCN        bool
	)

	cmd := &cobra.Command{
		Use:     "rm [name]",
		Aliases: []string{"remove"},
		Short:   "Remove Docker client certificate",
		Long:    removeDescription,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runRemove(args, tsaToken, tsaPassword, byCN)
		},
	}

	flags := cmd.Flags()
	flags.StringVarP(&tsaToken, "token", "t", "", "Token")
	flags.StringVarP(&tsaPassword, "password", "p", "", "Password")
	flags.BoolVarP(&byCN, "by-cn", "c", false, "Revoke certificate by Common Name instead of serial number")

	return cmd
}

func runRemove(args []string, tsaToken, tsaPassword string, byCN bool) error {
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

	if byCN {
		if err := clt.CertRevokeByCN(token, crt.CN); err != nil {
			if !isAlreadyRevokedOrExpired(err) {
				return err
			}
			fmt.Fprintf(os.Stderr, "Warning: certificate already revoked or expired on server, removing locally\n")
		}
	} else {
		certificate, err := pkix.NewCertificateFromPEMFile(cfg.TLS.CrtFile)
		if err != nil {
			return err
		}

		if err := clt.RevokeCertificate(token, int(certificate.Crt.SerialNumber.Int64())); err != nil {
			if !isAlreadyRevokedOrExpired(err) {
				return err
			}
			fmt.Fprintf(os.Stderr, "Warning: certificate already revoked or expired on server, removing locally\n")
		}
	}

	if err := s.RemoveCert(args[0]); err != nil {
		return err
	}

	return os.RemoveAll(cfg.Profile.CertDir)
}

// isAlreadyRevokedOrExpired checks if the TSA error indicates the certificate
// is no longer in a valid state (already revoked or expired).
func isAlreadyRevokedOrExpired(err error) bool {
	msg := strings.ToLower(err.Error())
	return strings.Contains(msg, "invalid revocation status") ||
		strings.Contains(msg, "already revoked")
}

var removeDescription = `
Remove Docker client certificate

`
