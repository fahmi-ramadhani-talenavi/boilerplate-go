package fileutil

import (
	"io"

	"github.com/xuri/excelize/v2"
)

// GenerateXLSX writes data to XLSX format
func GenerateXLSX(w io.Writer, sheetName string, data [][]string) error {
	f := excelize.NewFile()
	defer f.Close()

	index, err := f.NewSheet(sheetName)
	if err != nil {
		return err
	}

	for i, row := range data {
		for j, col := range row {
			cell, _ := excelize.CoordinatesToCellName(j+1, i+1)
			f.SetCellValue(sheetName, cell, col)
		}
	}

	f.SetActiveSheet(index)
	return f.Write(w)
}
