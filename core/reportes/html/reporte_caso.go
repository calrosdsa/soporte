package html

import (
	"bytes"
	"html/template"
	"log"
	"soporte-go/core/model/caso"
	"soporte-go/core/model/user"
	"soporte-go/core/model/ws"
)



type CasoReporte struct {
	Caso         caso.Caso
	UsuariosCaso []user.UserForList
	Messages []ws.Message
}

func HtmlCasoReporte(buf *bytes.Buffer, c caso.Caso, u []user.UserForList,m []ws.Message) {
	// users := []User{
	// 	{Name: "John Smith", Age: 35, Email: "john@example.com"},
	// 	{Name: "Jane Doe", Age: 27, Email: "jane@example.com"},
	// 	{Name: "Bob Johnson", Age: 42, Email: "bob@example.com"},
	// }
	// tmpl, err := template.ParseFiles("template.html")
	tmpl, err := template.ParseFiles("/home/rootuser/soporte/app/template.html")

	if err != nil {
		log.Println(tmpl)
	}
	data := CasoReporte{
		Caso:         c,
		UsuariosCaso: u,
		Messages: m,
	}

	// execute the HTML template with the data
	err = tmpl.Execute(buf, data)
	if err != nil {
		log.Println("ERROR ------------", err)
	}

	// if err != nil {
	// http.Error(w, err.Error(), http.StatusInternalServerError)
	// return
	// }
}
