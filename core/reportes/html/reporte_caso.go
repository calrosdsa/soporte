package html

import (
	"bytes"
	"html/template"
	"log"
	"soporte-go/core/model/caso"
)

type User struct {
	Name  string
	Age   int
	Email string
}

func HtmlCasoReporte(buf *bytes.Buffer, c caso.Caso) {
	// users := []User{
	// 	{Name: "John Smith", Age: 35, Email: "john@example.com"},
	// 	{Name: "Jane Doe", Age: 27, Email: "jane@example.com"},
	// 	{Name: "Bob Johnson", Age: 42, Email: "bob@example.com"},
	// }
	tmpl, err := template.ParseFiles("template.html")
	if err != nil {
		log.Println(tmpl)
	}
	// execute the HTML template with the data
	err = tmpl.Execute(buf, c)
	if err != nil{
		log.Println(err)
	}

	// if err != nil {
	// http.Error(w, err.Error(), http.StatusInternalServerError)
	// return
	// }
}
