package cert

import (
	"errors"
	"fmt"

	"github.com/juliengk/go-cert/pkix"
	"github.com/kassisol/tsa/pkg/adf"
	"github.com/kassisol/twic/pkg/date"
	"github.com/kassisol/twic/pkg/format"
	"github.com/kassisol/twic/storage"
	"github.com/spf13/cobra"
)

func newListCommand() *cobra.Command {
	var outputFormat string

	cmd := &cobra.Command{
		Use:     "ls",
		Aliases: []string{"list"},
		Short:   "List Docker client certificates",
		Long:    listDescription,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runList(args, outputFormat)
		},
	}

	cmd.Flags().StringVarP(&outputFormat, "format", "f", "table", "Output format: table or json")

	return cmd
}

type certListEntry struct {
	Name   string `json:"name"`
	Type   string `json:"type"`
	CN     string `json:"cn"`
	TSAURL string `json:"tsa_url"`
	Expire string `json:"expire,omitempty"`
}

func runList(args []string, outputFormat string) error {
	if len(args) > 0 {
		return errors.New("this command takes no arguments")
	}

	if outputFormat != "table" && outputFormat != "json" {
		return fmt.Errorf("unsupported format %q: use table or json", outputFormat)
	}

	cfg := adf.NewClient()
	if err := cfg.Init(); err != nil {
		return err
	}

	s, err := storage.NewDriver("sqlite", cfg.App.Dir.Root)
	if err != nil {
		return err
	}
	defer s.End()

	certs, err := s.ListCerts()
	if err != nil {
		return err
	}

	var entries []certListEntry
	for _, c := range certs {
		var expire string

		cfg.SetName(c.Name)

		certificate, err := pkix.NewCertificateFromPEMFile(cfg.TLS.CrtFile)
		if err == nil {
			expire = date.ExpireDateString(certificate.Crt.NotAfter)
		}

		entries = append(entries, certListEntry{
			Name:   c.Name,
			Type:   c.Type,
			CN:     c.CN,
			TSAURL: c.TSAURL,
			Expire: expire,
		})
	}

	if outputFormat == "json" {
		return format.PrintJSON(entries)
	}

	if len(entries) > 0 {
		var rows [][]string
		for _, e := range entries {
			rows = append(rows, []string{e.Name, e.Type, e.CN, e.TSAURL, e.Expire})
		}
		format.Table([]string{"NAME", "TYPE", "CN", "TSA URL", "EXPIRE"}, rows)
	}

	return nil
}

var listDescription = `
List Docker client certificates

`
