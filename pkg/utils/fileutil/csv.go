package fileutil

import (
	"encoding/csv"
	"io"
)

// GenerateCSV writes data to CSV format
func GenerateCSV(w io.Writer, data [][]string) error {
	writer := csv.NewWriter(w)
	defer writer.Flush()

	for _, row := range data {
		if err := writer.Write(row); err != nil {
			return err
		}
	}
	return nil
}
