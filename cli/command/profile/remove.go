package profile

import (
	"errors"

	"github.com/kassisol/tsa/pkg/adf"
	"github.com/kassisol/twic/storage"
	"github.com/spf13/cobra"
)

func newRemoveCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "rm [name]",
		Aliases: []string{"remove"},
		Short:   "Remove Docker profile",
		Long:    removeDescription,
		RunE:    runRemove,
	}

	return cmd
}

func runRemove(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return errors.New("this command requires exactly one argument")
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

	return s.RemoveProfile(args[0])
}

var removeDescription = `
Remove Docker profile

`
