package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/jung-kurt/gofpdf"
)

func main() {
	// Create a new PDF document
	pdf := gofpdf.New("P", "mm", "A4", "")

	// Set document properties
	pdf.SetTitle("Sample PDF Report",true)
	pdf.SetAuthor("John Doe",true)
	pdf.SetSubject("Sales Report",true)
	header := func() {
		// Set font and position for header text
		pdf.SetFont("Arial", "B", 14)
		pdf.CellFormat(0, 10, "Sales Report", "", 0, "C", false, 0, "")

		// Add logo
		// pdf.Image("logo.png", 165, 8, 30, 0, false, "", 0, "")
		pdf.Image("logo.png",10,10,10,10,false,"",0,"")

		// Draw a horizontal line below the header
		pdf.Ln(20)
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
	headers := []string{"Date", "Customer", "Product", "Quantity", "Price", "Total"}
	pdf.SetFont("Arial", "B", 12)
	pdf.Ln(20)
	for _, header := range headers {
		pdf.CellFormat(31.5, 7, header, "1", 0, "C", false, 0, "")
	}
	pdf.Ln(-1)

	// Add table data
	data := [][]string{
		{"2023-03-01", "John Doe", "Product A", "5", "$10.00", "$50.00"},
		{"2023-03-02", "Jane Doe", "Product B", "3", "$31.5.00", "$60.00"},
		{"2023-03-03", "Bob Smith", "Product C", "2", "$15.00", "$31.5.00"},
	}
	pdf.SetFont("Arial", "", 10)
	for _, row := range data {
		for _, datum := range row {
			pdf.CellFormat(31.5, 7, datum, "1", 0, "", false, 0, "")
		}
		pdf.Ln(-1)
	}
	pdf.Text(100,100,"Hello")

	pdf.AddPage()

	resp, err := http.Get("https://www.thewowstyle.com/wp-content/uploads/2015/01/nature-images.jpg")
    if err != nil {
        fmt.Println(err)
        return
    }
    defer resp.Body.Close()

    // Read the image data
    imageBytes, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        fmt.Println(err)
        return
    }

    // Register the image in the PDF document
    pdf.RegisterImageReader("Example Image","jpg" , bytes.NewReader(imageBytes))
	pdf.Image("Example Image", 10, 20, 50, 0, false, "", 0, "")

    // Output the PDF document to a file
    err = pdf.OutputFileAndClose("example.pdf")
    if err != nil {
        fmt.Println(err)
        return
    }



	// pdf.Image("logo.png",float64(10),float64(10),float64(10),float64(10),false,"",int(0),"")

	// Save PDF document to file
	// err := pdf.OutputFileAndClose("report.pdf")
	// if err != nil {
	// 	fmt.Println("Error creating PDF file:", err)
	// }
}
