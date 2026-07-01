package controller

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/stackriv/dev-tools/internal/business/model"
	"github.com/stackriv/dev-tools/internal/config"
	"github.com/stackriv/dev-tools/internal/pkg"
)

func Logs(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		if r.URL.Path != "/logs" {
			err := pkg.ErrorMessage(http.StatusNotFound)
			config.RenderTemplate(w, "error", model.PageData{Error: model.ErrorData{Code: err["code"], Message: err["msg"]}})
			fmt.Println(http.StatusNotFound, err["msg"])
			return
		}

		config.RenderTemplate(w, "logs", model.PageData{
			Title:       "Log Viewer",
			Description: "View and filter logs from your Docker containers.",
			Page:        "logs",
		})
	}
}

func GetContainers(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/api/containers" {
		err := pkg.ErrorMessage(http.StatusNotFound)
		config.RenderTemplate(w, "error", model.PageData{Error: model.ErrorData{Code: err["code"], Message: err["msg"]}})
		fmt.Println(http.StatusNotFound, err["msg"])
		return
	}

	w.Header().Set("Content-Type", "application/json")

	cmd := exec.Command("docker", "ps", "-a", "--format", "{{.ID}}|{{.Names}}|{{.Image}}|{{.State}}")
	out, err := cmd.Output()
	if err != nil {
		err1 := pkg.ErrorMessage(http.StatusInternalServerError)
		config.RenderTemplate(w, "error", model.PageData{Error: model.ErrorData{Code: err1["code"], Message: err.Error()}})
		fmt.Println(http.StatusInternalServerError, err.Error())
		return
	}

	var containers []model.Container
	scanner := bufio.NewScanner(strings.NewReader(string(out)))
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.SplitN(line, "|", 4)
		if len(parts) == 4 {
			containers = append(containers, model.Container{
				ID:    parts[0],
				Name:  parts[1],
				Image: parts[2],
				State: parts[3],
			})
		}
	}

	if containers == nil {
		containers = []model.Container{}
	}

	err = json.NewEncoder(w).Encode(map[string]interface{}{
		"containers": containers,
	})
	if err != nil {
		err1 := pkg.ErrorMessage(http.StatusInternalServerError)
		config.RenderTemplate(w, "error", model.PageData{Error: model.ErrorData{Code: err1["code"], Message: err1["msg"]}})
		fmt.Println(err1)
	}
}

func GetLogs(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/api/logs" {
		err := pkg.ErrorMessage(http.StatusNotFound)
		config.RenderTemplate(w, "error", model.PageData{Error: model.ErrorData{Code: err["code"], Message: err["msg"]}})
		fmt.Println(http.StatusNotFound, err["msg"])
		return
	}

	w.Header().Set("Content-Type", "application/json")

	var body struct {
		Container string `json:"container"`
		Lines     int    `json:"lines"`
		Search    string `json:"search"`
		Since     string `json:"since"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.Container == "" {
		err1 := pkg.ErrorMessage(http.StatusBadRequest)
		config.RenderTemplate(w, "error", model.PageData{Error: model.ErrorData{Code: err1["code"], Message: "container required"}})
		fmt.Println(http.StatusBadRequest, "container required")
		return
	}

	if body.Lines == 0 {
		body.Lines = 100
	}

	args := []string{"logs", "--tail", strconv.Itoa(body.Lines), "--timestamps"}
	if body.Since != "" {
		args = append(args, "--since", body.Since)
	}
	args = append(args, body.Container)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "docker", args...)
	out, _ := cmd.CombinedOutput()

	lines := strings.Split(string(out), "\n")
	var filtered []string

	for _, line := range lines {
		if line == "" {
			continue
		}
		if body.Search != "" && !strings.Contains(strings.ToLower(line), strings.ToLower(body.Search)) {
			continue
		}
		filtered = append(filtered, line)
	}

	if filtered == nil {
		filtered = []string{}
	}

	err := json.NewEncoder(w).Encode(map[string]interface{}{
		"logs":  filtered,
		"count": len(filtered),
	})
	if err != nil {
		err1 := pkg.ErrorMessage(http.StatusInternalServerError)
		config.RenderTemplate(w, "error", model.PageData{Error: model.ErrorData{Code: err1["code"], Message: err1["msg"]}})
		fmt.Println(err1)
	}
}
