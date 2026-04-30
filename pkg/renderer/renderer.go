package renderer

import (
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"
)

func RenderTemplate(w http.ResponseWriter, templateFile string) {
	tmpl, err := templateFor(templateFile)
	if err != nil {
		fmt.Println("Something went wrong: ", err)
		tmpl, _ = templateFor("error")
	}

	err = tmpl.Execute(w, nil)
	if err != nil {
		fmt.Println("Something went wrong: ", err)
		w.Write([]byte("Oops, something went wrong..."))
	}
}

var cachedTemplates = map[string]*template.Template{}

func templateFor(pageName string) (*template.Template, error) {
	tmpl, cached := cachedTemplates[pageName]
	if cached {
		fmt.Println("Using cache: ", pageName)
	} else {
		srcFiles := []string{fmt.Sprintf("templates/%s.page.tmpl", pageName)}
		layoutFiles, err := filepath.Glob("templates/*.layout.tmpl")
		if err != nil {
			return nil, err
		}
		srcFiles = append(srcFiles, layoutFiles...)
		tmpl, err = template.ParseFiles(srcFiles...)
		if err != nil {
			return nil, err
		}

		cachedTemplates[pageName] = tmpl
		fmt.Println("Cache template: ", pageName)
	}

	return tmpl, nil
}
