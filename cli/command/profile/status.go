package profile

import (
	"errors"
	"fmt"
	"os"

	"github.com/kassisol/twic/pkg/format"
	"github.com/spf13/cobra"
)

func newStatusCommand() *cobra.Command {
	var outputFormat string

	cmd := &cobra.Command{
		Use:   "status",
		Short: "Display Docker environment variables if set",
		Long:  statusDescription,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runStatus(args, outputFormat)
		},
	}

	cmd.Flags().StringVarP(&outputFormat, "format", "f", "table", "Output format: table or json")

	return cmd
}

type dockerStatus struct {
	Host      string `json:"docker_host,omitempty"`
	TLSVerify string `json:"docker_tls_verify,omitempty"`
	CertPath  string `json:"docker_cert_path,omitempty"`
	Set       bool   `json:"set"`
}

func runStatus(args []string, outputFormat string) error {
	if len(args) > 0 {
		return errors.New("this command takes no arguments")
	}

	if outputFormat != "table" && outputFormat != "json" {
		return fmt.Errorf("unsupported format %q: use table or json", outputFormat)
	}

	host := os.Getenv("DOCKER_HOST")
	tlsVerify := os.Getenv("DOCKER_TLS_VERIFY")
	certPath := os.Getenv("DOCKER_CERT_PATH")

	isSet := len(host) > 0 && len(tlsVerify) > 0 && len(certPath) > 0

	if outputFormat == "json" {
		entry := dockerStatus{Set: isSet}
		if isSet {
			entry.Host = host
			entry.TLSVerify = tlsVerify
			entry.CertPath = certPath
		}
		return format.PrintJSON(entry)
	}

	if isSet {
		fmt.Printf("DOCKER_HOST=%s\n", host)
		fmt.Printf("DOCKER_TLS_VERIFY=%s\n", tlsVerify)
		fmt.Printf("DOCKER_CERT_PATH=%s\n", certPath)
	} else {
		fmt.Println("Docker variables are not set")
	}

	return nil
}

var statusDescription = `
Display Docker environment variables if set

`
