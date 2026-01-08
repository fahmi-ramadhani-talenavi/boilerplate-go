package fileutil

import (
	"io"

	"github.com/go-pdf/fpdf"
)

// GeneratePDF writes data to PDF format
func GeneratePDF(w io.Writer, title string, data [][]string) error {
	pdf := fpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(40, 10, title)
	pdf.Ln(12)

	pdf.SetFont("Arial", "", 12)
	for _, row := range data {
		for _, col := range row {
			pdf.Cell(40, 10, col)
		}
		pdf.Ln(8)
	}

	return pdf.Output(w)
}
