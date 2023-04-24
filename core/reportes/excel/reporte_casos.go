package excel

import (
	"fmt"
	"log"
	"soporte-go/core/model"
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
	
    CreateSheet(casos2,"Sheet1",f)
    CreateSheet(casos,"Sheet2",f)
	
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
	var (
		cliente string
		update interface{}
		fechaInicio interface{}
		// fechaFin interface{}

	)
	
	f.NewSheet(sheet)
	f.SetColWidth(sheet, "A", "A", 10)
	f.SetColWidth(sheet, "B", "B", 55)
	f.SetColWidth(sheet, "C", "C", 27)
	f.SetColWidth(sheet, "D", "F", 17)
	// f.SetColWidth(sheet, "E", "F", 15)
	f.SetColWidth(sheet, "G", "G", 27)
	f.SetColWidth(sheet, "H", "H", 22)
	f.SetColWidth(sheet, "I", "I", 17)
	// f.SetColWidth(sheet, "G", "J", 20)


	titleStyle, err := f.NewStyle(&excelize.Style{
        Font:      &excelize.Font{Color: "1f7f3b", Bold: true, Family: "Arial"},
        Fill:      excelize.Fill{Type: "pattern", Color: []string{"E6F4EA"}, Pattern: 1},
        Alignment: &excelize.Alignment{Vertical: "center", Horizontal: "center"},
        Border:    []excelize.Border{{Type: "top", Style: 2, Color: "1f7f3b"},{Type: "left", Style: 2, Color: "1f7f3b"},
		{Type: "bottom", Style: 2, Color: "1f7f3b"},{Type: "right", Style: 2, Color: "1f7f3b"}},
        // Border:    []excelize.Border{{Type: "Bottom", Style: 2, Color: "1f7f3b"}},
    });
	if err != nil {
		log.Println(err)
	}
    // set style for the 'SUNDAY' to 'SATURDAY'
    if err := f.SetCellStyle(sheet, "A2", "I2", titleStyle); err != nil {
        log.Println(err)
        return
    }
	headers := []string{"Key","Asunto","Usuario asignado","F. Creacion","Ult. Actualizacion","F. Inicio","Author","Proyecto","Estado"}
	cell, err := excelize.CoordinatesToCellName(1, 2)
	if err != nil{
		log.Println(err)
	}
	f.SetSheetRow(sheet, cell, &headers)
	
	if err != nil {
		fmt.Println(err)
	}
	for idx,c := range casos {
		cliente = *c.ClienteName + " "+ *c.ClienteApellido
		var funcionario string
		if c.FuncionarioName != nil {
			funcionario = *c.FuncionarioName + " " + *c.FuncionarioApellido
		}
		if c.UpdatedOn == nil { update = "" }else{ update = *c.UpdatedOn }
		if c.FechaInicio == nil { fechaInicio = "" }else{ fechaInicio = *c.FechaInicio }
		// if c.FechaFin == nil { fechaFin = "" }else{ fechaFin = *c.FechaFin }

		slice := []interface{}{c.Key,c.Titulo,funcionario,c.CreatedOn,update,fechaInicio,cliente,*c.ProyectoName,
			GetCasoEstado(*c.Estado)}
		cell, err := excelize.CoordinatesToCellName(1, idx+3)
		// f.SetColWidth("Sheet1","B", 35)
		if err != nil {
			log.Panicln(err)
		}
		f.SetSheetRow(sheet, cell, &slice)
	}
}

func GetCasoEstado(estado int) (res string){
	switch estado {
	case int(model.NoResuelto):
		return "No Resuelto"
	case int(model.Resuelto):
		return "Resuelto"
	case int(model.Pendiente):
		return "No iniciado"
	default:
		return "En curso"
	}
}