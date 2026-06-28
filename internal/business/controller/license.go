package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/stackriv/dev-tools/internal/business/model"
	"github.com/stackriv/dev-tools/internal/config"
	"github.com/stackriv/dev-tools/internal/globale"
	"github.com/stackriv/dev-tools/internal/pkg"
)

func License(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		if r.URL.Path != "/license" {
			err := pkg.ErrorMessage(http.StatusNotFound)
			config.RenderTemplate(w, "error", model.PageData{Error: model.ErrorData{Code: err["code"], Message: err["msg"]}})
			fmt.Println(http.StatusNotFound, err["msg"])
			return
		}

		config.RenderTemplate(w, "license", model.PageData{
			Title:       "License Generator",
			Description: "Generate an open source license for your project.",
			Page:        "license",
		})
	}
}

func GenerateLicense(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/api/license" {
		err := pkg.ErrorMessage(http.StatusNotFound)
		config.RenderTemplate(w, "error", model.PageData{Error: model.ErrorData{Code: err["code"], Message: err["msg"]}})
		fmt.Println(http.StatusNotFound, err["msg"])
		return
	}

	w.Header().Set("Content-Type", "application/json")

	var body struct {
		Type    string `json:"type"`
		Author  string `json:"author"`
		Year    string `json:"year"`
		Project string `json:"project"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		err := pkg.ErrorMessage(http.StatusBadRequest)
		config.RenderTemplate(w, "error", model.PageData{Error: model.ErrorData{Code: err["code"], Message: err["msg"]}})
		fmt.Println(http.StatusBadRequest, err["msg"])
		return
	}

	tmpl, ok := globale.LicenseTemplates[body.Type]
	if !ok {
		tmpl = globale.LicenseTemplates["mit"]
	}

	author := body.Author
	if strings.TrimSpace(author) == "" {
		author = "Stackriv"
	}

	content := fmt.Sprintf(tmpl, body.Year, author)

	err := json.NewEncoder(w).Encode(map[string]string{"content": content})
	if err != nil {
		err := pkg.ErrorMessage(http.StatusInternalServerError)
		config.RenderTemplate(w, "error", model.PageData{Error: model.ErrorData{Code: err["code"], Message: err["msg"]}})
		fmt.Println(err)
	}
}
