package config

import (
	"bytes"
	"net/http"
	"path/filepath"
	"text/template"

	"github.com/stackriv/go-web-starter/internal/business/model"
	"github.com/stackriv/go-web-starter/internal/pkg"
)

// Global variable to hold the application configuration
var appConfig *Config

// NewAppConfig sets the global application configuration
func NewAppConfig(app *Config) {
	appConfig = app
}

// RenderTemplate renders a template with the given name and data
func RenderTemplate(w http.ResponseWriter, tmplName string, tmplData interface{}) {
	templateCache := appConfig.TemplateCache
	tmpl, ok := templateCache[tmplName+".page.tmpl"]

	if !ok {
		err := pkg.ErrorMessage(http.StatusInternalServerError)
		RenderTemplate(w, "error", model.Starter{Error: model.ErrorData{Code: err["code"], Message: err["msg"]}})
		return
	}

	buffer := new(bytes.Buffer)
	err := tmpl.Execute(buffer, tmplData)
	if err != nil {
		err := pkg.ErrorMessage(http.StatusInternalServerError)
		RenderTemplate(w, "error", model.Starter{Error: model.ErrorData{Code: err["code"], Message: err["msg"]}})
		return
	}
	_, err = buffer.WriteTo(w)
	if err != nil {
		err := pkg.ErrorMessage(http.StatusInternalServerError)
		RenderTemplate(w, "error", model.Starter{Error: model.ErrorData{Code: err["code"]}})
		return
	}
}

// CreateTemplateCache creates a template cache from the provided template files
func CreateTemplateCache() (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}
	pages, err := filepath.Glob("./templates/*.page.tmpl")
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		name := filepath.Base(page)
		tmpl := template.Must(template.ParseFiles(page))

		layouts, err := filepath.Glob("./templates/layouts/*.layout.tmpl")
		if err != nil {
			return nil, err
		}

		if len(layouts) > 0 {
			_, err := tmpl.ParseGlob("./templates/layouts/*.layout.tmpl")
			if err != nil {
				return nil, err
			}
		}

		cache[name] = tmpl
	}

	return cache, nil
}
