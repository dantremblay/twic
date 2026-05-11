package access

import (
	"errors"
	"fmt"

	"github.com/kassisol/twic/pkg/input"
	"github.com/kassisol/tsa/client"
	"github.com/spf13/cobra"
)

func NewCommand() *cobra.Command {
	var (
		tsaURL      string
		tsaTTL      int
		tsaUsername string
		tsaPassword string
	)

	cmd := &cobra.Command{
		Use:   "access",
		Short: "Get TSA access token",
		Long:  accessDescription,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runAccess(args, tsaURL, tsaTTL, tsaUsername, tsaPassword)
		},
	}

	flags := cmd.Flags()
	flags.StringVarP(&tsaURL, "tsa-url", "c", "", "TSA URL")
	flags.IntVarP(&tsaTTL, "ttl", "t", 1440, "Token TTL")
	flags.StringVarP(&tsaUsername, "username", "u", "", "Username")
	flags.StringVarP(&tsaPassword, "password", "p", "", "Password")

	return cmd
}

func runAccess(args []string, tsaURL string, tsaTTL int, tsaUsername, tsaPassword string) error {
	if len(args) > 0 {
		return errors.New("this command takes no arguments")
	}

	tsaurl := tsaURL
	if len(tsaurl) == 0 {
		tsaurl = input.ReadInput("TSA URL")
	}

	username := tsaUsername
	if len(username) == 0 {
		username = input.ReadInput("Username")
	}

	password := tsaPassword
	if len(password) == 0 {
		password = input.ReadPassword("Password")
	}

	if len(username) == 0 {
		return errors.New("empty username is not allowed")
	}

	if len(password) == 0 {
		return errors.New("empty password is not allowed")
	}

	clt, err := client.New(tsaurl)
	if err != nil {
		return err
	}

	if err := clt.GetDirectory(); err != nil {
		return err
	}

	token, err := clt.GetToken(username, password, tsaTTL)
	if err != nil {
		return err
	}

	fmt.Println(token)
	return nil
}

var accessDescription = `
The **twic access** command has subcommands for Getting TSA access token.

To see help for a subcommand, use:

    twic access [command] --help

For full details on using twic access visit Harbormaster's online documentation.

`
