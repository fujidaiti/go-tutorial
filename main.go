package main

import (
	"fmt"
	"html/template"
	"net/http"
)

func main() {
	http.HandleFunc("/", home)
	http.HandleFunc("/about", about)
	http.ListenAndServe(":8080", nil)
}

func home(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "home.page.tmpl")
}

func about(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "about.page.tmpl")
}

func renderTemplate(w http.ResponseWriter, templateFile string) {
	template, _ := template.ParseFiles("template/" + templateFile)
	err := template.Execute(w, nil)
	if err != nil {
		fmt.Println("Something went wrong: ", err)
		errorTemplate, _ := template.ParseFiles("template/error.page.tmpl")
		errorTemplate.Execute(w, nil)
	}
}
