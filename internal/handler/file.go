package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/user/go-boilerplate/internal/dto"
	"github.com/user/go-boilerplate/pkg/apperror"
	"github.com/user/go-boilerplate/pkg/utils/fileutil"
)

type FileHandler struct{}

func NewFileHandler() *FileHandler {
	return &FileHandler{}
}

func (h *FileHandler) RegisterRoutes(g *echo.Group) {
	g.GET("/export/csv", h.ExportCSV)
	g.GET("/export/xlsx", h.ExportXLSX)
	g.GET("/export/pdf", h.ExportPDF)
}

func (h *FileHandler) ExportCSV(c echo.Context) error {
	data := getSampleData()

	c.Response().Header().Set(echo.HeaderContentDisposition, "attachment; filename=export.csv")
	c.Response().Header().Set(echo.HeaderContentType, "text/csv")

	if err := fileutil.GenerateCSV(c.Response().Writer, data); err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error: dto.ErrorDetail{
				Code:    string(apperror.ErrCodeInternal),
				Message: "Failed to generate CSV",
			},
		})
	}

	return nil
}

func (h *FileHandler) ExportXLSX(c echo.Context) error {
	data := getSampleData()

	c.Response().Header().Set(echo.HeaderContentDisposition, "attachment; filename=export.xlsx")
	c.Response().Header().Set(echo.HeaderContentType, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")

	if err := fileutil.GenerateXLSX(c.Response().Writer, "Export", data); err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error: dto.ErrorDetail{
				Code:    string(apperror.ErrCodeInternal),
				Message: "Failed to generate XLSX",
			},
		})
	}

	return nil
}

func (h *FileHandler) ExportPDF(c echo.Context) error {
	data := getSampleData()

	c.Response().Header().Set(echo.HeaderContentDisposition, "attachment; filename=export.pdf")
	c.Response().Header().Set(echo.HeaderContentType, "application/pdf")

	if err := fileutil.GeneratePDF(c.Response().Writer, "User Export", data); err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error: dto.ErrorDetail{
				Code:    string(apperror.ErrCodeInternal),
				Message: "Failed to generate PDF",
			},
		})
	}

	return nil
}

func getSampleData() [][]string {
	return [][]string{
		{"ID", "Name", "Email"},
		{"1", "John Doe", "john@example.com"},
		{"2", "Jane Smith", "jane@example.com"},
	}
}
