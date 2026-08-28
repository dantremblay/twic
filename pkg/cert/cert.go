package cert

import (
	"crypto/x509/pkix"
	"fmt"
	"slices"
	"strings"
)

// Subject holds the certificate subject fields extracted from a CA certificate.
type Subject struct {
	Country            string
	Province           string
	Locality           string
	Organization       string
	OrganizationalUnit string
}

// SubjectFromCA safely extracts the required subject fields from a CA
// certificate subject, returning an error if any mandatory field is missing
// instead of panicking on an out-of-range slice index.
func SubjectFromCA(name pkix.Name) (Subject, error) {
	fields := []struct {
		label  string
		values []string
	}{
		{"Country", name.Country},
		{"Province", name.Province},
		{"Locality", name.Locality},
		{"Organization", name.Organization},
		{"OrganizationalUnit", name.OrganizationalUnit},
	}

	for _, f := range fields {
		if len(f.values) == 0 {
			return Subject{}, fmt.Errorf("CA certificate subject is missing %s", f.label)
		}
	}

	return Subject{
		Country:            name.Country[0],
		Province:           name.Province[0],
		Locality:           name.Locality[0],
		Organization:       name.Organization[0],
		OrganizationalUnit: GetOU(name.OrganizationalUnit[0]),
	}, nil
}

func GetOU(ou string) string {
	words := []string{
		"Certificate",
		"Authority",
	}

	oldou := strings.Split(ou, " ")

	if len(oldou) > 1 {
		var newou []string

		for _, word := range oldou {
			if !slices.Contains(words, word) {
				newou = append(newou, word)
			}
		}

		if len(newou) > 0 {
			return strings.Join(newou, " ")
		}
	}

	return ou
}
