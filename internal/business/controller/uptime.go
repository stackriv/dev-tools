package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/stackriv/dev-tools/internal/business/model"
	"github.com/stackriv/dev-tools/internal/config"
	"github.com/stackriv/dev-tools/internal/pkg"
)

func Uptime(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		if r.URL.Path != "/uptime" {
			err := pkg.ErrorMessage(http.StatusNotFound)
			config.RenderTemplate(w, "error", model.PageData{Error: model.ErrorData{Code: err["code"], Message: err["msg"]}})
			fmt.Println(http.StatusNotFound, err["msg"])
			return
		}

		config.RenderTemplate(w, "uptime", model.PageData{
			Title:       "Uptime Monitor",
			Description: "Check the availability and response time of your services.",
			Page:        "uptime",
		})
	}
}

func CheckUptime(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/api/uptime" {
		err := pkg.ErrorMessage(http.StatusNotFound)
		config.RenderTemplate(w, "error", model.PageData{Error: model.ErrorData{Code: err["code"], Message: err["msg"]}})
		fmt.Println(http.StatusNotFound, err["msg"])
		return
	}

	w.Header().Set("Content-Type", "application/json")

	var body struct {
		Targets []struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"targets"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil || len(body.Targets) == 0 {
		err1 := pkg.ErrorMessage(http.StatusBadRequest)
		config.RenderTemplate(w, "error", model.PageData{Error: model.ErrorData{Code: err1["code"], Message: err.Error()}})
		fmt.Println(http.StatusBadRequest, err.Error())
		return
	}

	client := &http.Client{
		Timeout: 10 * time.Second,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return nil
		},
	}

	results := make([]model.UptimeTarget, len(body.Targets))
	var wg sync.WaitGroup

	for i, target := range body.Targets {
		wg.Add(1)
		go func(idx int, name, url string) {
			defer wg.Done()
			result := model.UptimeTarget{Name: name, URL: url}

			start := time.Now()
			resp, err := client.Get(url)
			elapsed := time.Since(start).Milliseconds()

			result.Ms = elapsed

			if err != nil {
				result.Status = "down"
				result.Error = err.Error()
			} else {
				resp.Body.Close()
				result.Code = resp.StatusCode
				if resp.StatusCode >= 200 && resp.StatusCode < 400 {
					result.Status = "up"
				} else {
					result.Status = "degraded"
				}
			}

			results[idx] = result
		}(i, target.Name, target.URL)
	}

	wg.Wait()

	err := json.NewEncoder(w).Encode(map[string]interface{}{
		"results":    results,
		"checked_at": time.Now().UTC().Format("2006-01-02 15:04:05 UTC"),
	})
	if err != nil {
		err1 := pkg.ErrorMessage(http.StatusInternalServerError)
		config.RenderTemplate(w, "error", model.PageData{Error: model.ErrorData{Code: err1["code"], Message: err1["msg"]}})
		fmt.Println(err1)
	}
}
