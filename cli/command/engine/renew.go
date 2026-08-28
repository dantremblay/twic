package engine

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/juliengk/go-cert/helpers"
	"github.com/juliengk/go-cert/pkix"
	"github.com/kassisol/tsa/client"
	"github.com/kassisol/tsa/pkg/adf"
	"github.com/kassisol/twic/pkg/cert"
	"github.com/kassisol/twic/pkg/input"
	"github.com/kassisol/twic/pkg/sysutil"
	"github.com/kassisol/twic/pkg/urlutil"
	"github.com/spf13/cobra"
)

func newRenewCommand() *cobra.Command {
	var (
		duration    int
		tsaToken    string
		tsaUsername string
		tsaPassword string
	)

	cmd := &cobra.Command{
		Use:   "renew",
		Short: "Renew Docker engine certificate",
		Long:  renewDescription,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runRenew(args, duration, tsaToken, tsaUsername, tsaPassword)
		},
	}

	flags := cmd.Flags()
	flags.IntVarP(&duration, "duration", "d", 12, "Certificate duration (in months)")
	flags.StringVarP(&tsaToken, "token", "t", "", "Token")
	flags.StringVarP(&tsaUsername, "username", "u", "", "Username")
	flags.StringVarP(&tsaPassword, "password", "p", "", "Password")

	return cmd
}

func runRenew(args []string, duration int, tsaToken, tsaUsername, tsaPassword string) error {
	if !sysutil.IsRoot() {
		return errors.New("you must be root to run engine subcommand")
	}

	if len(args) > 0 {
		return errors.New("this command takes no arguments")
	}

	certtype := "engine"

	var username, password string
	if len(tsaToken) == 0 {
		username = tsaUsername
		if len(username) == 0 {
			username = input.ReadInput("Username")
		}

		password = tsaPassword
		if len(password) == 0 {
			password = input.ReadPassword("Password")
		}

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

	// Load the existing engine certificate to recover its identity.
	oldcert, err := pkix.NewCertificateFromPEMFile(cfg.TLS.CrtFile)
	if err != nil {
		return err
	}

	// Recover the TSA URL from the certificate's CRL distribution point,
	// the same way engine info/rm do.
	if len(oldcert.Crt.CRLDistributionPoints) == 0 {
		return errors.New("cannot determine TSA URL: certificate has no CRL distribution point")
	}

	crldp := oldcert.Crt.CRLDistributionPoints[0]

	url, err := urlutil.Parse(crldp)
	if err != nil {
		return err
	}

	tsaurl := fmt.Sprintf("%s://%s", url.Scheme, url.Host)
	if url.Port != 443 {
		tsaurl = fmt.Sprintf("%s://%s:%d", url.Scheme, url.Host, url.Port)
	}

	// Extract subject fields and preserve the Subject Alternative Names.
	subject, err := cert.SubjectFromCA(oldcert.Crt.Subject)
	if err != nil {
		return err
	}

	certcn := oldcert.Crt.Subject.CommonName
	altnames := sanStrings(oldcert)

	clt, err := client.New(tsaurl)
	if err != nil {
		return err
	}

	if err := clt.GetDirectory(); err != nil {
		return err
	}

	token := tsaToken
	if len(token) == 0 {
		token, err = clt.GetToken(username, password, 5)
		if err != nil {
			return err
		}
	}

	// Generate a fresh key and CSR carrying the original identity forward.
	key, err := pkix.NewKey(4096)
	if err != nil {
		return err
	}

	keyBytes, err := key.ToPEM()
	if err != nil {
		return err
	}

	csr, err := helpers.CreateCSR(subject.Country, subject.Province, subject.Locality, subject.Organization, subject.OrganizationalUnit, certcn, "", altnames, key)
	if err != nil {
		return err
	}

	// The TSA server refuses to issue a certificate while a valid one with
	// the same distinguished name still exists, so the old certificate must
	// be revoked before requesting the new one. Tolerate a certificate that
	// is already revoked or expired.
	if err := clt.RevokeCertificate(token, int(oldcert.Crt.SerialNumber.Int64())); err != nil {
		if !isAlreadyRevokedOrExpired(err) {
			return err
		}
	}

	crt, err := clt.GetCertificate(token, certtype, csr.Bytes, duration)
	if err != nil {
		return fmt.Errorf("the old certificate was revoked but issuing a new one failed: %w\nRe-run this command or 'twic engine create' to obtain a new certificate", err)
	}

	// Persist the new key and certificate to temporary files first, then
	// atomically move them into place so a partial write cannot break the
	// daemon's TLS configuration.
	newKeyFile := cfg.TLS.KeyFile + ".new"
	newCrtFile := cfg.TLS.CrtFile + ".new"

	defer func() {
		os.Remove(newKeyFile)
		os.Remove(newCrtFile)
	}()

	if err := pkix.ToPEMFile(newKeyFile, keyBytes, 0400); err != nil {
		return err
	}

	if err := pkix.ToPEMFile(newCrtFile, []byte(crt), 0444); err != nil {
		return err
	}

	if err := os.Rename(newKeyFile, cfg.TLS.KeyFile); err != nil {
		return err
	}

	if err := os.Rename(newCrtFile, cfg.TLS.CrtFile); err != nil {
		return err
	}

	fmt.Println("Docker engine certificate renewed in the directory", cfg.CertsDir, ".")
	fmt.Println("\nRestart the Docker daemon for the new certificate to take effect (e.g. systemctl restart docker).")

	return nil
}

// sanStrings returns the Subject Alternative Names (DNS names and IP
// addresses) of a certificate as a slice of strings suitable for CreateCSR.
func sanStrings(certificate *pkix.Certificate) []string {
	var altnames []string

	altnames = append(altnames, certificate.Crt.DNSNames...)

	for _, ip := range certificate.Crt.IPAddresses {
		altnames = append(altnames, ip.String())
	}

	return altnames
}

// isAlreadyRevokedOrExpired checks if the TSA error indicates the certificate
// is no longer in a valid state (already revoked or expired).
func isAlreadyRevokedOrExpired(err error) bool {
	msg := strings.ToLower(err.Error())
	return strings.Contains(msg, "invalid revocation status") ||
		strings.Contains(msg, "already revoked")
}

var renewDescription = `
Renew Docker engine certificate

`
