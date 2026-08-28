package main

import (
	"fmt"
	"os"

	"github.com/kassisol/twic/cli/command/commands"
	"github.com/spf13/cobra"
)

func main() {
	scPath := os.Getenv("TWIC_SC_PATH")
	if scPath == "" {
		scPath = "dist/shellcompletion"
	}
	bashTarget := fmt.Sprintf("%s/bash", scPath)

	if err := os.MkdirAll(scPath, 0755); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	cmd := &cobra.Command{Use: "twic"}
	commands.AddCommands(cmd)
	cmd.DisableAutoGenTag = true

	if err := cmd.GenBashCompletionFile(bashTarget); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
