package renderer

import (
	"fmt"
	"html/template"
	"net/http"
)

func RenderTemplate(w http.ResponseWriter, templateFile string) {
	template, _ := template.ParseFiles("templates/" + templateFile)
	err := template.Execute(w, nil)
	if err != nil {
		fmt.Println("Something went wrong: ", err)
		errorTemplate, _ := template.ParseFiles("templates/error.page.tmpl")
		errorTemplate.Execute(w, nil)
	}
}
