package controller

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/stackriv/dev-tools/internal/business/model"
	"github.com/stackriv/dev-tools/internal/business/service"
	"github.com/stackriv/dev-tools/internal/config"
	"github.com/stackriv/dev-tools/internal/pkg"
)

func Markdown(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		if r.URL.Path != "/markdown" {
			err := pkg.ErrorMessage(http.StatusNotFound)
			config.RenderTemplate(w, "error", model.PageData{Error: model.ErrorData{Code: err["code"], Message: err["msg"]}})
			fmt.Println(http.StatusNotFound, err["msg"])
			return
		}

		config.RenderTemplate(w, "markdown", model.PageData{
			Title:       "Markdown Editor",
			Description: "Write Markdown and preview the result in real time.",
			Page:        "markdown",
		})
	}
}

func RenderMarkdown(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/api/markdown" {
		err := pkg.ErrorMessage(http.StatusNotFound)
		config.RenderTemplate(w, "error", model.PageData{Error: model.ErrorData{Code: err["code"], Message: err["msg"]}})
		fmt.Println(http.StatusNotFound, err["msg"])
		return
	}

	w.Header().Set("Content-Type", "application/json")

	var body struct {
		Content string `json:"content"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		err1 := pkg.ErrorMessage(http.StatusBadRequest)
		config.RenderTemplate(w, "error", model.PageData{Error: model.ErrorData{Code: err1["code"], Message: "invalid request"}})
		fmt.Println(http.StatusBadRequest, "invalid request")
		return
	}

	html := service.MarkdownToHTML(body.Content)
	err := json.NewEncoder(w).Encode(map[string]string{"html": html})
	if err != nil {
		err1 := pkg.ErrorMessage(http.StatusInternalServerError)
		config.RenderTemplate(w, "error", model.PageData{Error: model.ErrorData{Code: err1["code"], Message: err1["msg"]}})
		fmt.Println(err1)
	}
}
