package format

import (
	"encoding/json"
	"fmt"
	"os"
	"text/tabwriter"
)

// PrintJSON marshals data as indented JSON to stdout.
func PrintJSON(data any) error {
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	return enc.Encode(data)
}

// Table writes rows as tab-aligned columns to stdout.
// headers is the header row, rows is a slice of string slices.
func Table(headers []string, rows [][]string) {
	w := tabwriter.NewWriter(os.Stdout, 20, 1, 2, ' ', 0)

	for i, h := range headers {
		if i > 0 {
			fmt.Fprint(w, "\t")
		}
		fmt.Fprint(w, h)
	}
	fmt.Fprintln(w)

	for _, row := range rows {
		for i, col := range row {
			if i > 0 {
				fmt.Fprint(w, "\t")
			}
			fmt.Fprint(w, col)
		}
		fmt.Fprintln(w)
	}

	w.Flush()
}
