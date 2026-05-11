package main

import (
	"fmt"
	"os"

	"github.com/kassisol/twic/cli/command/commands"
	"github.com/spf13/cobra"
)

func main() {
	scPath := "/tmp/twic/shellcompletion"
	bashTarget := fmt.Sprintf("%s/bash", scPath)

	if err := os.MkdirAll(scPath, 0755); err != nil {
		fmt.Println(err)
	}

	cmd := &cobra.Command{Use: "twic"}
	commands.AddCommands(cmd)
	cmd.DisableAutoGenTag = true

	if err := cmd.GenBashCompletionFile(bashTarget); err != nil {
		fmt.Println(err)
	}
}
