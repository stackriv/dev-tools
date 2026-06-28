package controller

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/stackriv/dev-tools/internal/business/model"
	"github.com/stackriv/dev-tools/internal/config"
	"github.com/stackriv/dev-tools/internal/pkg"
)

func UUID(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		if r.URL.Path != "/uuid" {
			err := pkg.ErrorMessage(http.StatusNotFound)
			config.RenderTemplate(w, "error", model.PageData{Error: model.ErrorData{Code: err["code"], Message: err["msg"]}})
			fmt.Println(http.StatusNotFound, err["msg"])
			return
		}

		config.RenderTemplate(w, "uuid", model.PageData{
			Title:       "UUID Generator",
			Description: "Generate random UUIDs (v4).",
			Page:        "uuid",
		})
	}
}

func GenerateUUID(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/api/uuid" {
		err := pkg.ErrorMessage(http.StatusNotFound)
		config.RenderTemplate(w, "error", model.PageData{Error: model.ErrorData{Code: err["code"], Message: err["msg"]}})
		fmt.Println(http.StatusNotFound, err["msg"])
		return
	}

	w.Header().Set("Content-Type", "application/json")

	countStr := r.URL.Query().Get("count")
	count, err := strconv.Atoi(countStr)
	if err != nil || count < 1 || count > 100 {
		count = 1
	}

	uuids := make([]string, count)
	for i := range uuids {
		uuids[i] = newUUID()
	}

	err = json.NewEncoder(w).Encode(map[string]interface{}{"uuids": uuids})
	if err != nil {
		err := pkg.ErrorMessage(http.StatusInternalServerError)
		config.RenderTemplate(w, "error", model.PageData{Error: model.ErrorData{Code: err["code"], Message: err["msg"]}})
		fmt.Println(err)
	}
}

func newUUID() string {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		return ""
	}

	b[6] = (b[6] & 0x0f) | 0x40 // version 4
	b[8] = (b[8] & 0xbf) | 0x80 // variant

	return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
}
