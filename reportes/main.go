package main

import (
	"html/template"
	"net/http"
)

type User struct {
	Name    string
	Age     int
	Email   string
}

func main() {
	// create a new HTTP server
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// define the data to be passed to the HTML template
		users := []User{
			{Name: "John Smith", Age: 35, Email: "john@example.com"},
			{Name: "Jane Doe", Age: 27, Email: "jane@example.com"},
			{Name: "Bob Johnson", Age: 42, Email: "bob@example.com"},
		}

		// parse the HTML template file
		tmpl, err := template.ParseFiles("template.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// execute the HTML template with the data
		err = tmpl.Execute(w, users)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})

	// start the HTTP server
	http.ListenAndServe(":8080", nil)
}