package urlutil

import (
	"net/url"
	"strconv"
	"strings"
)

// URL is a parsed URL with the port extracted as an integer.
type URL struct {
	Scheme string
	Host   string
	Port   int
	Path   string
}

// Parse parses a raw URL string and extracts the host and port separately.
func Parse(rawurl string) (*URL, error) {
	u, err := url.Parse(rawurl)
	if err != nil {
		return nil, err
	}

	host := u.Hostname()
	port := 80

	if p := u.Port(); p != "" {
		port, _ = strconv.Atoi(p)
	} else if u.Scheme == "https" {
		port = 443
	}

	// Fallback: if Hostname() is empty, try splitting Host manually
	if host == "" && u.Host != "" {
		parts := strings.SplitN(u.Host, ":", 2)
		host = parts[0]
	}

	return &URL{
		Scheme: u.Scheme,
		Host:   host,
		Port:   port,
		Path:   u.Path,
	}, nil
}
