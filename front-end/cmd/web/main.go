package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
)

// main starts the web server and listens on port 80.
func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		render(w, "test.page.gohtml")
	})

	fmt.Println("Starting fron end service on port 80")

	err := http.ListenAndServe(":80", nil)
	if err != nil {
		log.Panic(err)
	}
}

// render writes the given template to w, along with the base layout and any partials.
func render(w http.ResponseWriter, t string) {
	partials := []string{
		"./cmd/web/templates/base.layout.gohtml",
		"./cmd/web/templates/header.partial.gohtml",
		"./cmd/web/templates/footer.partial.gohtml",
	}

	// Create list of templates
	var templateSlice []string

	// Add the given file to the templates
	templateSlice = append(templateSlice, "./cmd/web/templates/"+t)

	// Add the partials to the templates
	templateSlice = append(templateSlice, partials...)

	// Parse the files into a template object
	tmpl, err := template.ParseFiles(templateSlice...)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Execute the template object
	if err := tmpl.Execute(w, nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
