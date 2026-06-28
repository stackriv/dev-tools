package controller

import (
	"fmt"
	"net/http"

	"github.com/stackriv/go-web-starter/internal/business/model"
	"github.com/stackriv/go-web-starter/internal/config"
	"github.com/stackriv/go-web-starter/internal/pkg"
)

func Starter(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		if r.URL.Path != "/" {
			err := pkg.ErrorMessage(http.StatusNotFound)
			config.RenderTemplate(w, "error", model.Starter{Error: model.ErrorData{Code: err["code"], Message: err["msg"]}})
			fmt.Println(http.StatusNotFound, err["msg"])
			return
		}

		names := make(map[string]string)
		names["owner"] = "Stackriv"
		config.RenderTemplate(w, "starter", model.Starter{StringData: names})
		fmt.Println(http.StatusOK, "OK")
	}
}
