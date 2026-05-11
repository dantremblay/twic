package profile

import (
	"errors"
	"fmt"

	"github.com/kassisol/tsa/pkg/adf"
	"github.com/kassisol/twic/pkg/format"
	"github.com/kassisol/twic/storage"
	"github.com/spf13/cobra"
)

func newListCommand() *cobra.Command {
	var outputFormat string

	cmd := &cobra.Command{
		Use:     "ls",
		Aliases: []string{"list"},
		Short:   "List Docker profiles",
		Long:    listDescription,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runList(args, outputFormat)
		},
	}

	cmd.Flags().StringVarP(&outputFormat, "format", "f", "table", "Output format: table or json")

	return cmd
}

type profileListEntry struct {
	Name       string `json:"name"`
	CertName   string `json:"cert_name"`
	DockerHost string `json:"docker_host"`
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

	profiles, err := s.ListProfiles()
	if err != nil {
		return err
	}

	var entries []profileListEntry
	for _, p := range profiles {
		entries = append(entries, profileListEntry{
			Name:       p.Name,
			CertName:   p.Cert.Name,
			DockerHost: p.DockerHost,
		})
	}

	if outputFormat == "json" {
		return format.PrintJSON(entries)
	}

	if len(entries) > 0 {
		var rows [][]string
		for _, e := range entries {
			rows = append(rows, []string{e.Name, e.CertName, e.DockerHost})
		}
		format.Table([]string{"NAME", "CERTIFICATE NAME", "DOCKER HOST"}, rows)
	}

	return nil
}

var listDescription = `
List Docker profiles

`
