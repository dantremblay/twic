package engine

import (
	"errors"

	"github.com/kassisol/twic/pkg/sysutil"
	"github.com/spf13/cobra"
)

func newRenewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "renew",
		Short: "Renew Docker engine certificate",
		Long:  renewDescription,
		RunE:  runRenew,
	}

	return cmd
}

func runRenew(cmd *cobra.Command, args []string) error {
	if !sysutil.IsRoot() {
		return errors.New("you must be root to run engine subcommand")
	}

	if len(args) > 0 {
		return errors.New("this command takes no arguments")
	}

	return errors.New("not implemented yet")
}

var renewDescription = `
Renew Docker engine certificate

`
