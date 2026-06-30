package controller

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"

	"github.com/stackriv/dev-tools/internal/business/model"
	"github.com/stackriv/dev-tools/internal/config"
	"github.com/stackriv/dev-tools/internal/pkg"
)

func DNS(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		if r.URL.Path != "/dns" {
			err := pkg.ErrorMessage(http.StatusNotFound)
			config.RenderTemplate(w, "error", model.PageData{Error: model.ErrorData{Code: err["code"], Message: err["msg"]}})
			fmt.Println(http.StatusNotFound, err["msg"])
			return
		}

		config.RenderTemplate(w, "dns", model.PageData{
			Title:       "DNS Lookup",
			Description: "Resolve DNS records for any domain.",
			Page:        "dns",
		})
	}
}

func LookupDNS(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/api/dns" {
		err := pkg.ErrorMessage(http.StatusNotFound)
		config.RenderTemplate(w, "error", model.PageData{Error: model.ErrorData{Code: err["code"], Message: err["msg"]}})
		fmt.Println(http.StatusNotFound, err["msg"])
		return
	}

	w.Header().Set("Content-Type", "application/json")

	var body struct {
		Domain string   `json:"domain"`
		Types  []string `json:"types"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.Domain == "" {
		err := pkg.ErrorMessage(http.StatusBadRequest)
		config.RenderTemplate(w, "error", model.PageData{Error: model.ErrorData{Code: err["code"], Message: "domain required"}})
		fmt.Println(err)
		return
	}

	if len(body.Types) == 0 {
		body.Types = []string{"A", "AAAA", "MX", "TXT", "NS", "CNAME"}
	}

	results := make(map[string]model.DNSResult)

	for _, t := range body.Types {
		result := model.DNSResult{Type: t}

		switch t {
		case "A":
			ips, err := net.LookupHost(body.Domain)
			if err != nil {
				result.Error = err.Error()
			} else {
				for _, ip := range ips {
					if net.ParseIP(ip).To4() != nil {
						result.Records = append(result.Records, ip)
					}
				}
			}

		case "AAAA":
			ips, err := net.LookupHost(body.Domain)
			if err != nil {
				result.Error = err.Error()
			} else {
				for _, ip := range ips {
					if net.ParseIP(ip).To4() == nil {
						result.Records = append(result.Records, ip)
					}
				}
			}

		case "MX":
			mxs, err := net.LookupMX(body.Domain)
			if err != nil {
				result.Error = err.Error()
			} else {
				for _, mx := range mxs {
					result.Records = append(result.Records, mx.Host)
				}
			}

		case "TXT":
			txts, err := net.LookupTXT(body.Domain)
			if err != nil {
				result.Error = err.Error()
			} else {
				result.Records = txts
			}

		case "NS":
			nss, err := net.LookupNS(body.Domain)
			if err != nil {
				result.Error = err.Error()
			} else {
				for _, ns := range nss {
					result.Records = append(result.Records, ns.Host)
				}
			}

		case "CNAME":
			cname, err := net.LookupCNAME(body.Domain)
			if err != nil {
				result.Error = err.Error()
			} else {
				result.Records = []string{cname}
			}
		}

		if result.Records == nil {
			result.Records = []string{}
		}
		results[t] = result
	}

	err := json.NewEncoder(w).Encode(map[string]interface{}{
		"domain":  body.Domain,
		"results": results,
	})
	if err != nil {
		err := pkg.ErrorMessage(http.StatusInternalServerError)
		config.RenderTemplate(w, "error", model.PageData{Error: model.ErrorData{Code: err["code"], Message: err["msg"]}})
		fmt.Println(err)
	}
}
