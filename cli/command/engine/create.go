package engine

import (
	"errors"
	"fmt"
	"os"
	"slices"
	"strings"

	"github.com/juliengk/go-cert/helpers"
	"github.com/juliengk/go-cert/pkix"
	"github.com/kassisol/tsa/client"
	"github.com/kassisol/tsa/pkg/adf"
	"github.com/kassisol/twic/pkg/cert"
	"github.com/kassisol/twic/pkg/input"
	"github.com/kassisol/twic/pkg/sysutil"
	"github.com/spf13/cobra"
)

func newCreateCommand() *cobra.Command {
	var (
		certCN       string
		certAltNames string
		duration     int
		tsaURL       string
		tsaToken     string
		tsaUsername   string
		tsaPassword   string
	)

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create Docker engine certificate",
		Long:  createDescription,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runCreate(args, certCN, certAltNames, duration, tsaURL, tsaToken, tsaUsername, tsaPassword)
		},
	}

	flags := cmd.Flags()
	flags.StringVarP(&certCN, "common-name", "n", "", "Certificate Common Name")
	flags.StringVarP(&certAltNames, "alt-names", "a", "", "Certificate Alternative Names")
	flags.IntVarP(&duration, "duration", "d", 12, "Certificate duration (in months)")
	flags.StringVarP(&tsaURL, "tsa-url", "c", "", "TSA URL")
	flags.StringVarP(&tsaToken, "token", "t", "", "Token")
	flags.StringVarP(&tsaUsername, "username", "u", "", "Username")
	flags.StringVarP(&tsaPassword, "password", "p", "", "Password")

	return cmd
}

func runCreate(args []string, certCN, certAltNames string, duration int, tsaURL, tsaToken, tsaUsername, tsaPassword string) (retErr error) {
	if !sysutil.IsRoot() {
		return errors.New("you must be root to run engine subcommand")
	}

	if len(args) > 0 {
		return errors.New("this command takes no arguments")
	}

	certtype := "engine"

	certcn := certCN
	if len(certcn) == 0 {
		certcn = input.ReadInput("Common Name (CN)")
	}

	certaltnames := certAltNames
	if len(certaltnames) == 0 {
		certaltnames = input.ReadInput("Alt Names")
	}

	tsaurl := tsaURL
	if len(tsaurl) == 0 {
		tsaurl = input.ReadInput("TSA URL")
	}

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
	}

	cfg := adf.NewEngine()
	if err := cfg.Init(); err != nil {
		return err
	}

	if len(tsaToken) == 0 {
		if len(username) == 0 {
			return errors.New("empty username is not allowed")
		}
		if len(password) == 0 {
			return errors.New("empty password is not allowed")
		}
	}

	// Cleanup on failure
	defer func() {
		if retErr != nil {
			if sysutil.FileExists(cfg.TLS.CaFile) {
				os.Remove(cfg.TLS.CaFile)
			}
			if sysutil.FileExists(cfg.TLS.KeyFile) {
				os.Remove(cfg.TLS.KeyFile)
			}
		}
	}()

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

	caCrt, err := clt.GetCACertificate()
	if err != nil {
		return err
	}

	if err := pkix.ToPEMFile(cfg.TLS.CaFile, []byte(caCrt), 0444); err != nil {
		return err
	}

	key, err := helpers.CreateKey(4096, cfg.TLS.KeyFile)
	if err != nil {
		return err
	}

	caCertificate, err := pkix.NewCertificateFromPEM([]byte(caCrt))
	if err != nil {
		return err
	}

	ca := caCertificate.Crt.Subject
	ou := cert.GetOU(ca.OrganizationalUnit[0])

	ans := strings.Split(certaltnames, ",")
	if !slices.Contains(ans, certcn) {
		ans = append(ans, certcn)
	}

	csr, err := helpers.CreateCSR(ca.Country[0], ca.Province[0], ca.Locality[0], ca.Organization[0], ou, certcn, "", ans, key)
	if err != nil {
		return err
	}

	crt, err := clt.GetCertificate(token, certtype, csr.Bytes, duration)
	if err != nil {
		return err
	}

	if err := pkix.ToPEMFile(cfg.TLS.CrtFile, []byte(crt), 0444); err != nil {
		return err
	}

	fmt.Println("Docker engine certificates created in the directory", cfg.CertsDir, ".")
	fmt.Printf("\nTo configure the Docker Daemon, add the following parameters:\n\n--tlsverify --tlscacert=%s --tlskey=%s --tlscert=%s -H tcp://0.0.0.0:2376\n\n", cfg.TLS.CaFile, cfg.TLS.KeyFile, cfg.TLS.CrtFile)

	return nil
}

var createDescription = `
Create Docker engine certificate

`
