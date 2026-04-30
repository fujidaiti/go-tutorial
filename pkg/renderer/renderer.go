package renderer

import (
	"fmt"
	"html/template"
	"net/http"
)

func RenderTemplate(w http.ResponseWriter, templateFile string) {
	tmpl, err := templateFor(templateFile)
	if err != nil {
		fmt.Println("Something went wrong: ", err)
		tmpl, _ = templateFor("templates/error.page.tmpl")
	}

	err = tmpl.Execute(w, nil)
	if err != nil {
		fmt.Println("Something went wrong: ", err)
		w.Write([]byte("Oops, something went wrong..."))
	}
}

var cache = make(map[string]*template.Template)

func templateFor(templateFile string) (*template.Template, error) {
	tmpl, cached := cache[templateFile]
	if !cached {
		files := []string{
			fmt.Sprintf("templates/%s", templateFile),
			"templates/base.layout.tmpl",
		}

		var err error
		tmpl, err = template.ParseFiles(files...)
		if err != nil {
			return nil, err
		}

		cache[templateFile] = tmpl
		fmt.Println("Cache template: ", templateFile)
	}

	return tmpl, nil
}
