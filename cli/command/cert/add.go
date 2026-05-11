package cert

import (
	"errors"
	"fmt"
	"net/url"
	"os"

	"github.com/juliengk/go-cert/helpers"
	"github.com/juliengk/go-cert/pkix"
	"github.com/kassisol/tsa/client"
	"github.com/kassisol/tsa/pkg/adf"
	"github.com/kassisol/twic/pkg/cert"
	"github.com/kassisol/twic/pkg/input"
	"github.com/kassisol/twic/storage"
	"github.com/kassisol/twic/pkg/sysutil"
	"github.com/kassisol/twic/pkg/validate"
	"github.com/spf13/cobra"
)

func newAddCommand() *cobra.Command {
	var (
		tsaURL      string
		tsaToken    string
		tsaUsername string
		tsaPassword string
	)

	cmd := &cobra.Command{
		Use:   "add [name]",
		Short: "Add Docker client certificate",
		Long:  addDescription,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runAdd(args, tsaURL, tsaToken, tsaUsername, tsaPassword)
		},
	}

	flags := cmd.Flags()
	flags.StringVarP(&tsaURL, "tsa-url", "c", "", "TSA URL")
	flags.StringVarP(&tsaToken, "token", "t", "", "Token")
	flags.StringVarP(&tsaUsername, "username", "u", "", "Username")
	flags.StringVarP(&tsaPassword, "password", "p", "", "Password")

	return cmd
}

func runAdd(args []string, tsaURL, tsaToken, tsaUsername, tsaPassword string) error {
	if sysutil.IsRoot() {
		return errors.New("you must not be root to add a client certificate type")
	}

	if len(args) != 1 {
		return errors.New("this command requires exactly one argument")
	}

	name := args[0]
	certtype := "client"

	tsaurl := tsaURL
	if len(tsaurl) == 0 {
		tsaurl = input.ReadInput("TSA URL")
	}

	username := tsaUsername
	if len(username) == 0 {
		username = input.ReadInput("Username")
	}

	var password string
	if len(tsaToken) == 0 {
		password = tsaPassword
		if len(password) == 0 {
			password = input.ReadPassword("Password")
		}
	}

	certcn := username

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

	if err := validate.Name(name); err != nil {
		return err
	}

	if existing, _ := s.GetCert(name); existing.Name != "" {
		return fmt.Errorf("name %q already exists", name)
	}

	tu, err := url.Parse(tsaurl)
	if err != nil {
		return err
	}
	if tu.Scheme != "https" {
		return errors.New("TSA URL scheme should be https")
	}

	if len(username) == 0 {
		return errors.New("empty username is not allowed")
	}

	if len(tsaToken) == 0 && len(password) == 0 {
		return errors.New("empty password is not allowed")
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

	csr, err := helpers.CreateCSR(ca.Country[0], ca.Province[0], ca.Locality[0], ca.Organization[0], ou, certcn, "", []string{}, key)
	if err != nil {
		return err
	}

	crt, err := clt.GetCertificate(token, certtype, csr.Bytes, 12)
	if err != nil {
		_ = os.RemoveAll(cfg.Profile.CertDir)
		return err
	}

	if err := pkix.ToPEMFile(cfg.TLS.CrtFile, []byte(crt), 0444); err != nil {
		return err
	}

	return s.AddCert(name, certtype, certcn, "", tsaurl)
}

var addDescription = `
Add Docker client certificate

`
