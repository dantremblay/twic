package cert

import (
	"slices"
	"strings"
)

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
