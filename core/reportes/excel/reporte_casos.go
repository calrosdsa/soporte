package excel

import (
	"fmt"
	"log"
	"soporte-go/core/model/caso"

	"bytes"

	"github.com/xuri/excelize/v2"
)

func ReporteCasosExcel(casos []caso.Caso,casos2 []caso.Caso,buffer *bytes.Buffer)(err error) {
	f := excelize.NewFile()
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()
	
    CreateSheet(casos,"Sheet1",f)
    CreateSheet(casos2,"Sheet2",f)
	
	// for idx, row := range [][]interface{}{
	// 	{nil, "Apple", "Orange", "Pear"}, {"Small", 2, 3, 3},
	// 	{"Normal", 5, 2, 4}, {"Large", 6, 7, 8},
	// } {
	// 	cell, err := excelize.CoordinatesToCellName(1, idx+1)
	// 	if err != nil {
	// 		fmt.Println(err)
	// 	}
	// 	f.SetSheetRow("Sheet1", cell, &row)
	// }
	// if err := f.AddChart("Sheet1", "E1", &excelize.Chart{
	// 	Type: excelize.Col3DClustered,
	// 	Series: []excelize.ChartSeries{
	// 		{
	// 			Name:       "Sheet1!$A$2",
	// 			Categories: "Sheet1!$B$1:$D$1",
	// 			Values:     "Sheet1!$B$2:$D$2",
	// 		},
	// 		{
	// 			Name:       "Sheet1!$A$3",
	// 			Categories: "Sheet1!$B$1:$D$1",
	// 			Values:     "Sheet1!$B$3:$D$3",
	// 		},
	// 		{
	// 			Name:       "Sheet1!$A$4",
	// 			Categories: "Sheet1!$B$1:$D$1",
	// 			Values:     "Sheet1!$B$4:$D$4",
	// 		}},
	// 	Title: excelize.ChartTitle{
	// 		Name: "Fruit 3D Clustered Column Chart",
	// 	},
	// }); err != nil {
	// 	fmt.Println(err)
	// }
	// f.SetActiveSheet(index)

	err = f.Write(buffer)
	if err != nil {
		log.Println(err)
	}
    return
	// Save spreadsheet by the given path.
	// if err := f.SaveAs("Book1.xlsx"); err != nil {
	// 	fmt.Println(err)
	// }
}


func CreateSheet(casos []caso.Caso,sheet string,f *excelize.File){
	var cliente string
	f.NewSheet(sheet)
	f.SetColWidth(sheet, "A", "D", 35)
	f.SetColWidth(sheet, "D", "E", 25)
	for idx,c := range casos {
		cliente = *c.ClienteName + " "+ *c.ClienteApellido
		var funcionario string
		if c.FuncionarioName != nil {
			funcionario = *c.FuncionarioName + " " + *c.FuncionarioApellido
		}
		slice := []interface{}{c.Titulo,c.Id,c.CreatedOn,*c.Estado,cliente,funcionario}
		cell, err := excelize.CoordinatesToCellName(1, idx+1)
		// f.SetColWidth("Sheet1","B", 35)
		if err != nil {
			fmt.Println(err)
		}
		f.SetSheetRow(sheet, cell, &slice)
	}
}