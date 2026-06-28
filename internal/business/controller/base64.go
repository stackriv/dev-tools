package controller

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/stackriv/dev-tools/internal/business/model"
	"github.com/stackriv/dev-tools/internal/config"
	"github.com/stackriv/dev-tools/internal/pkg"
)

func Base64Handler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		if r.URL.Path != "/base64" {
			err := pkg.ErrorMessage(http.StatusNotFound)
			config.RenderTemplate(w, "error", model.PageData{Error: model.ErrorData{Code: err["code"], Message: err["msg"]}})
			fmt.Println(http.StatusNotFound, err["msg"])
			return
		}

		config.RenderTemplate(w, "uuid", model.PageData{
			Title:       "Base64 Encoder / Decoder",
			Description: "Encode or decode Base64 strings.",
			Page:        "base64",
		})
	}
}

func ProcessBase64(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/api/base64" {
		err := pkg.ErrorMessage(http.StatusNotFound)
		config.RenderTemplate(w, "error", model.PageData{Error: model.ErrorData{Code: err["code"], Message: err["msg"]}})
		fmt.Println(http.StatusNotFound, err["msg"])
		return
	}

	w.Header().Set("Content-Type", "application/json")

	var body struct {
		Input  string `json:"input"`
		Action string `json:"action"` // "encode" or "decode"
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.Input == "" {
		err := pkg.ErrorMessage(http.StatusBadRequest)
		config.RenderTemplate(w, "error", model.PageData{Error: model.ErrorData{Code: err["code"], Message: err["msg"]}})
		fmt.Println(err)
		return
	}

	var result string
	var encodeErr string

	switch body.Action {
	case "decode":
		decoded, err := base64.StdEncoding.DecodeString(body.Input)
		if err != nil {
			decoded, err = base64.URLEncoding.DecodeString(body.Input)
			if err != nil {
				encodeErr = "invalid base64 input"
			}
		}
		if strings.TrimSpace(encodeErr) == "" {
			result = string(decoded)
		}
	case "encode":
		result = base64.StdEncoding.EncodeToString([]byte(body.Input))
	default:
		encodeErr = "invalid action"
	}

	if encodeErr != "" {
		err := pkg.ErrorMessage(http.StatusBadRequest)
		config.RenderTemplate(w, "error", model.PageData{Error: model.ErrorData{Code: err["code"], Message: err["msg"]}})
		fmt.Println(err)
		return
	}

	err := json.NewEncoder(w).Encode(map[string]interface{}{"result": result})
	if err != nil {
		err := pkg.ErrorMessage(http.StatusInternalServerError)
		config.RenderTemplate(w, "error", model.PageData{Error: model.ErrorData{Code: err["code"], Message: err["msg"]}})
		fmt.Println(err)
	}
}
