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

func JWT(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		if r.URL.Path != "/jwt" {
			err := pkg.ErrorMessage(http.StatusNotFound)
			config.RenderTemplate(w, "error", model.PageData{Error: model.ErrorData{Code: err["code"], Message: err["msg"]}})
			fmt.Println(http.StatusNotFound, err["msg"])
			return
		}

		config.RenderTemplate(w, "jwt", model.PageData{
			Title:       "JWT Decoder",
			Description: "Decode and inspect a JWT token without sending it to any external service.",
			Page:        "jwt",
		})
	}
}

func DecodeJWT(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/api/jwt" {
		err := pkg.ErrorMessage(http.StatusNotFound)
		config.RenderTemplate(w, "error", model.PageData{Error: model.ErrorData{Code: err["code"], Message: err["msg"]}})
		fmt.Println(http.StatusNotFound, err["msg"])
		return
	}

	w.Header().Set("Content-Type", "application/json")

	var body struct {
		Token string `json:"token"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.Token == "" {
		err := pkg.ErrorMessage(http.StatusInternalServerError)
		config.RenderTemplate(w, "error", model.PageData{Error: model.ErrorData{Code: err["code"], Message: "token required"}})
		fmt.Println("token required")
		return
	}

	parts := strings.Split(strings.TrimSpace(body.Token), ".")
	if len(parts) != 3 {
		err := pkg.ErrorMessage(http.StatusBadRequest)
		config.RenderTemplate(w, "error", model.PageData{Error: model.ErrorData{Code: err["code"], Message: "invalid JWT format — expected 3 parts separated by dots"}})
		fmt.Println(http.StatusBadRequest, "invalid JWT format — expected 3 parts separated by dots")
		return
	}

	header, err := decodeSegment(parts[0])
	if err != nil {
		err1 := pkg.ErrorMessage(http.StatusBadRequest)
		config.RenderTemplate(w, "error", model.PageData{Error: model.ErrorData{Code: err1["code"], Message: "invalid header: " + err.Error()}})
		fmt.Println(http.StatusBadRequest, "invalid header: "+err.Error())
		return
	}

	payload, err := decodeSegment(parts[1])
	if err != nil {
		err1 := pkg.ErrorMessage(http.StatusBadRequest)
		config.RenderTemplate(w, "error", model.PageData{Error: model.ErrorData{Code: err1["code"], Message: "invalid payload: " + err.Error()}})
		fmt.Println(http.StatusBadRequest, "invalid payload: "+err.Error())
		return
	}

	// Print header and payload
	var headerMap, payloadMap interface{}
	_ = json.Unmarshal([]byte(header), &headerMap)
	_ = json.Unmarshal([]byte(payload), &payloadMap)

	headerPretty, _ := json.MarshalIndent(headerMap, "", "  ")
	payloadPretty, _ := json.MarshalIndent(payloadMap, "", "  ")

	err1 := json.NewEncoder(w).Encode(map[string]interface{}{
		"header":    string(headerPretty),
		"payload":   string(payloadPretty),
		"signature": parts[2],
		"valid":     true,
	})
	if err1 != nil {
		err := pkg.ErrorMessage(http.StatusInternalServerError)
		config.RenderTemplate(w, "error", model.PageData{Error: model.ErrorData{Code: err["code"], Message: err["msg"]}})
		fmt.Println(err)
	}
}

func decodeSegment(seg string) (string, error) {
	// Add padding if required
	switch len(seg) % 4 {
	case 2:
		seg += "=="
	case 3:
		seg += "="
	}
	decoded, err := base64.URLEncoding.DecodeString(seg)
	if err != nil {
		return "", err
	}
	return string(decoded), nil
}
