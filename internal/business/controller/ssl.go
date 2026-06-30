package controller

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/stackriv/dev-tools/internal/business/model"
	"github.com/stackriv/dev-tools/internal/config"
	"github.com/stackriv/dev-tools/internal/pkg"
)

func SSL(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		if r.URL.Path != "/ssl" {
			err := pkg.ErrorMessage(http.StatusNotFound)
			config.RenderTemplate(w, "error", model.PageData{Error: model.ErrorData{Code: err["code"], Message: err["msg"]}})
			fmt.Println(http.StatusNotFound, err["msg"])
			return
		}

		config.RenderTemplate(w, "ssl", model.PageData{
			Title:       "SSL Checker",
			Description: "Check SSL certificate validity and expiration for any domain.",
			Page:        "ssl",
		})
	}
}

func CheckSSL(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/api/ssl" {
		err := pkg.ErrorMessage(http.StatusNotFound)
		config.RenderTemplate(w, "error", model.PageData{Error: model.ErrorData{Code: err["code"], Message: err["msg"]}})
		fmt.Println(http.StatusNotFound, err["msg"])
		return
	}

	w.Header().Set("Content-Type", "application/json")

	var body struct {
		Domain string `json:"domain"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.Domain == "" {
		err := pkg.ErrorMessage(http.StatusBadRequest)
		config.RenderTemplate(w, "error", model.PageData{Error: model.ErrorData{Code: err["code"], Message: "domain required"}})
		fmt.Println(err)
		return
	}

	conn, err1 := tls.Dial("tcp", body.Domain+":443", &tls.Config{
		InsecureSkipVerify: false,
	})
	if err1 != nil {
		err := pkg.ErrorMessage(http.StatusBadRequest)
		config.RenderTemplate(w, "error", model.PageData{Error: model.ErrorData{Code: err["code"], Message: err1.Error()}})
		fmt.Println(err)
		return
	}
	defer conn.Close()

	cert := conn.ConnectionState().PeerCertificates[0]
	now := time.Now()
	daysLeft := int(cert.NotAfter.Sub(now).Hours() / 24)

	status := "valid"
	if daysLeft < 0 {
		status = "expired"
	} else if daysLeft <= 14 {
		status = "critical"
	} else if daysLeft <= 30 {
		status = "warning"
	}

	var sanList []string
	sanList = append(sanList, cert.DNSNames...)
	for _, ip := range cert.IPAddresses {
		sanList = append(sanList, ip.String())
	}

	err := json.NewEncoder(w).Encode(map[string]interface{}{
		"valid":       true,
		"status":      status,
		"domain":      body.Domain,
		"subject":     cert.Subject.CommonName,
		"issuer":      cert.Issuer.CommonName,
		"issued_at":   cert.NotBefore.Format("2006-01-02"),
		"expires_at":  cert.NotAfter.Format("2006-01-02"),
		"days_left":   daysLeft,
		"san":         sanList,
		"serial":      fmt.Sprintf("%X", cert.SerialNumber),
		"tls_version": tlsVersion(conn.ConnectionState().Version),
	})
	if err != nil {
		err := pkg.ErrorMessage(http.StatusInternalServerError)
		config.RenderTemplate(w, "error", model.PageData{Error: model.ErrorData{Code: err["code"], Message: err["msg"]}})
		fmt.Println(err)
	}
}

func tlsVersion(v uint16) string {
	switch v {
	case tls.VersionTLS10:
		return "TLS 1.0"
	case tls.VersionTLS11:
		return "TLS 1.1"
	case tls.VersionTLS12:
		return "TLS 1.2"
	case tls.VersionTLS13:
		return "TLS 1.3"
	default:
		return "Unknown"
	}
}
