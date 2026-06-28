package controller

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/stackriv/dev-tools/internal/business/model"
	"github.com/stackriv/dev-tools/internal/config"
	"github.com/stackriv/dev-tools/internal/pkg"
)

func Hash(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		if r.URL.Path != "/hash" {
			err := pkg.ErrorMessage(http.StatusNotFound)
			config.RenderTemplate(w, "error", model.PageData{Error: model.ErrorData{Code: err["code"], Message: err["msg"]}})
			fmt.Println(http.StatusNotFound, err["msg"])
			return
		}

		config.RenderTemplate(w, "uuid", model.PageData{
			Title:       "Hash Generator",
			Description: "Generate MD5, SHA1, SHA256 and SHA512 hashes.",
			Page:        "hash",
		})
	}
}

func GenerateHash(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/api/hash" {
		err := pkg.ErrorMessage(http.StatusNotFound)
		config.RenderTemplate(w, "error", model.PageData{Error: model.ErrorData{Code: err["code"], Message: err["msg"]}})
		fmt.Println(http.StatusNotFound, err["msg"])
		return
	}

	w.Header().Set("Content-Type", "application/json")

	var body struct {
		Input string `json:"input"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.Input == "" {
		err := pkg.ErrorMessage(http.StatusBadRequest)
		config.RenderTemplate(w, "error", model.PageData{Error: model.ErrorData{Code: err["code"], Message: "input required"}})
		fmt.Println(err)
		return
	}

	data := []byte(body.Input)

	md5sum := md5.Sum(data)
	sha1sum := sha1.Sum(data)
	sha256sum := sha256.Sum256(data)
	sha512sum := sha512.Sum512(data)

	err := json.NewEncoder(w).Encode(map[string]string{
		"md5":    fmt.Sprintf("%x", md5sum),
		"sha1":   fmt.Sprintf("%x", sha1sum),
		"sha256": fmt.Sprintf("%x", sha256sum),
		"sha512": fmt.Sprintf("%x", sha512sum),
	})
	if err != nil {
		err := pkg.ErrorMessage(http.StatusInternalServerError)
		config.RenderTemplate(w, "error", model.PageData{Error: model.ErrorData{Code: err["code"], Message: err["msg"]}})
		fmt.Println(err)
	}
}
