package validate

import (
	"errors"
	"regexp"
)

var nameRegex = regexp.MustCompile(`^[a-zA-Z0-9\-_]+$`)

// Name checks that a name contains only alphanumeric characters, hyphens, and underscores.
func Name(name string) error {
	if !nameRegex.MatchString(name) {
		return errors.New("name is not valid")
	}
	return nil
}

// Port checks that a port number is in the valid range.
func Port(port int) error {
	if port < 0 || port > 65535 {
		return errors.New("port is not valid")
	}
	return nil
}
