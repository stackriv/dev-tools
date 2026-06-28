package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/stackriv/dev-tools/internal/business/model"
	"github.com/stackriv/dev-tools/internal/config"
	"github.com/stackriv/dev-tools/internal/pkg"
)

func PackageJSON(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		if r.URL.Path != "/packagejson" {
			err := pkg.ErrorMessage(http.StatusNotFound)
			config.RenderTemplate(w, "error", model.PageData{Error: model.ErrorData{Code: err["code"], Message: err["msg"]}})
			fmt.Println(http.StatusNotFound, err["msg"])
			return
		}

		config.RenderTemplate(w, "packagejson", model.PageData{
			Title:       "package.json Generator",
			Description: "Generate a package.json file for your Node.js project.",
			Page:        "packagejson",
		})
	}
}

func GeneratePackageJSON(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/api/packagejson" {
		err := pkg.ErrorMessage(http.StatusNotFound)
		config.RenderTemplate(w, "error", model.PageData{Error: model.ErrorData{Code: err["code"], Message: err["msg"]}})
		fmt.Println(http.StatusNotFound, err["msg"])
		return
	}

	w.Header().Set("Content-Type", "application/json")

	var body struct {
		Name           string `json:"name"`
		Version        string `json:"version"`
		Description    string `json:"description"`
		Author         string `json:"author"`
		License        string `json:"license"`
		Type           string `json:"type"`
		IncludeScripts bool   `json:"includeScripts"`
		IncludeEngines bool   `json:"includeEngines"`
		IncludePrivate bool   `json:"includePrivate"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		err := pkg.ErrorMessage(http.StatusBadRequest)
		config.RenderTemplate(w, "error", model.PageData{Error: model.ErrorData{Code: err["code"], Message: err["msg"]}})
		fmt.Println(http.StatusBadRequest, err["msg"])
		return
	}

	pkg1 := map[string]interface{}{
		"name":        body.Name,
		"version":     body.Version,
		"description": body.Description,
		"author":      body.Author,
		"license":     body.License,
	}

	if strings.TrimSpace(body.Type) == "" {
		pkg1["type"] = body.Type
	}

	if body.IncludePrivate {
		pkg1["private"] = true
	}

	if body.IncludeScripts {
		pkg1["scripts"] = map[string]string{
			"start":  "node index.js",
			"dev":    "nodemon index.js",
			"build":  "tsc",
			"test":   "jest",
			"lint":   "eslint .",
			"format": "prettier --write .",
		}
	}

	if body.IncludeEngines {
		pkg1["engines"] = map[string]string{
			"node": ">=22.0.0",
			"npm":  ">=11.0.0",
		}
	}

	pkg1["keywords"] = []string{}
	pkg1["dependencies"] = map[string]string{}
	pkg1["devDependencies"] = map[string]string{}

	out, _ := json.MarshalIndent(pkg1, "", "  ")
	err := json.NewEncoder(w).Encode(map[string]string{"content": string(out)})
	if err != nil {
		err := pkg.ErrorMessage(http.StatusInternalServerError)
		config.RenderTemplate(w, "error", model.PageData{Error: model.ErrorData{Code: err["code"], Message: err["msg"]}})
		fmt.Println(err)
	}
}
