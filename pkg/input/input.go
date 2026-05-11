package input

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"golang.org/x/term"
)

// ReadInput prompts for and reads a line of text from stdin.
func ReadInput(label string) string {
	fmt.Printf("%s : ", label)

	reader := bufio.NewReader(os.Stdin)
	line, err := reader.ReadString('\n')
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	return strings.TrimRight(line, "\n")
}

// ReadPassword prompts for and reads a password from stdin with masked input.
func ReadPassword(label string) string {
	fmt.Printf("%s : ", label)

	pw, err := term.ReadPassword(int(os.Stdin.Fd()))
	fmt.Println() // newline after masked input
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	return string(pw)
}
