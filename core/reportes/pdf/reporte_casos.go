package pdf

import (
	"bytes"
	"fmt"
	"soporte-go/core/model/caso"

	"github.com/jung-kurt/gofpdf"
)

func ReporteCasos(casos []caso.Caso, buffer *bytes.Buffer) (err error) {
	pdf := gofpdf.New("P", "mm", "letter", "")

	// Set document properties
	pdf.SetTitle("PDF Report", true)
	pdf.SetAuthor("John Doe", true)
	pdf.SetSubject("Reporte Casos", true)
	header := func() {
		// Set font and position for header text
		pdf.SetFont("Arial", "B", 14)
		pdf.CellFormat(0, 10, "Sales Report", "", 0, "C", false, 0, "")
		// Add logo
		pdf.Image("logo.png", 10, 10, 10, 10, false, "", 0, "")

		// Draw a horizontal line below the header
		pdf.Ln(15)
		pdf.Line(10, 30, 200, 30)
	}
	pdf.SetHeaderFunc(header)

	footer := func() {
		// Set font and position for footer text
		pdf.SetY(-15)
		pdf.SetFont("Arial", "", 8)
		pdf.CellFormat(0, 10, "Page "+fmt.Sprintf("%d", pdf.PageNo()), "", 0, "C", false, 0, "")
	}
	pdf.SetFooterFunc(footer)
	// Add first page
	pdf.AddPage()
	// Add table headers
	headers := []string{"Asunto","Creado","Inicio", "Fin", "Cliente", "Funcionario"}
	pdf.SetFont("Arial", "B", 9)
	pdf.Ln(15)
	for _, header := range headers {
		if header == "Asunto" {
			pdf.CellFormat(40, 6, header, "1", 0, "C", false, 0, "")
		} else {
			pdf.CellFormat(30, 6, header, "1", 0, "C", false, 0, "")
		}
	}
	pdf.Ln(-1)

	// Add table data
	
	pdf.SetFont("Arial", "", 7)
	// pdf.SetTextColor(103,142,132)
	for _, caso := range casos {
		truncado := caso.Titulo
		if len(caso.Titulo) > 26 {
			truncado = caso.Titulo[:26] + "..."
		}
		pdf.SetTextColor(0, 0, 255)
		pdf.CellFormat(40, 6, truncado, "1", 0, "", false, 0, "https://google.com")
		pdf.SetTextColor(0, 0, 0)
		pdf.CellFormat(30, 6, caso.CreatedOn.Format("2006-01-02 15:04"), "1", 0, "", false, 0, "")
		if caso.FechaInicio == nil {
			pdf.CellFormat(30, 6, "No iniciado", "1", 0, "", false, 0, "")
		}else {
			pdf.CellFormat(30, 6, caso.FechaInicio.Format("2006-01-02 15:04"), "1", 0, "", false, 0, "")
		}
		if caso.FechaInicio == nil {
			pdf.CellFormat(30, 6, "No finalizado", "1", 0, "", false, 0, "")
		}else {
			pdf.CellFormat(30, 6, caso.FechaFin.Format("2006-01-02 15:04"), "1", 0, "", false, 0, "")
		}

		// pdf.CellFormat(30, 6, "", "1", 0, "", false, 0, "")
		pdf.CellFormat(30, 6, "*caso.ClienteId", "1", 0, "", false, 0, "")
		pdf.CellFormat(30, 6, "*caso.FuncionarioId", "1", 0, "", false, 0, "")
		pdf.Ln(-1)
	}
	pdf.Text(100, 100, "Hello")

	pdf.AddPage()
	// Output the PDF document to a file
	err = pdf.Output(buffer)
	if err != nil {
		fmt.Println(err)
		return
	}
	return
}
