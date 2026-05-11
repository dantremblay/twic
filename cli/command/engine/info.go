package engine

import (
	"errors"
	"fmt"

	"github.com/juliengk/go-cert/pkix"
	"github.com/kassisol/tsa/pkg/adf"
	"github.com/kassisol/twic/pkg/date"
	"github.com/kassisol/twic/pkg/format"
	"github.com/kassisol/twic/pkg/sysutil"
	"github.com/kassisol/twic/pkg/urlutil"
	"github.com/spf13/cobra"
)

func newInfoCommand() *cobra.Command {
	var outputFormat string

	cmd := &cobra.Command{
		Use:   "info",
		Short: "Information about Docker engine certificate",
		Long:  infoDescription,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runInfo(args, outputFormat)
		},
	}

	cmd.Flags().StringVarP(&outputFormat, "format", "f", "table", "Output format: table or json")

	return cmd
}

type engineInfoEntry struct {
	TSAURL string `json:"tsa_url"`
	CN     string `json:"cn"`
	Expire string `json:"expire"`
}

func runInfo(args []string, outputFormat string) error {
	if !sysutil.IsRoot() {
		return errors.New("you must be root to run engine subcommand")
	}

	if len(args) > 0 {
		return errors.New("this command takes no arguments")
	}

	if outputFormat != "table" && outputFormat != "json" {
		return fmt.Errorf("unsupported format %q: use table or json", outputFormat)
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

	cn := certificate.Crt.Subject.CommonName
	expire := date.ExpireDateString(certificate.Crt.NotAfter)

	entry := engineInfoEntry{
		TSAURL: tsaurl,
		CN:     cn,
		Expire: expire,
	}

	if outputFormat == "json" {
		return format.PrintJSON(entry)
	}

	fmt.Println("TSA URL:", entry.TSAURL)
	fmt.Println("CN:", entry.CN)
	fmt.Println("Expire:", entry.Expire)

	return nil
}

var infoDescription = `
Information about Docker engine certificate

`
