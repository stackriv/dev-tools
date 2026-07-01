package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/stackriv/dev-tools/internal/business/model"
	"github.com/stackriv/dev-tools/internal/config"
	"github.com/stackriv/dev-tools/internal/pkg"
)

func Readme(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		if r.URL.Path != "/readme" {
			err := pkg.ErrorMessage(http.StatusNotFound)
			config.RenderTemplate(w, "error", model.PageData{Error: model.ErrorData{Code: err["code"], Message: err["msg"]}})
			fmt.Println(http.StatusNotFound, err["msg"])
			return
		}

		config.RenderTemplate(w, "readme", model.PageData{
			Title:       "README Generator",
			Description: "Generate a professional README.md for your project.",
			Page:        "readme",
		})
	}
}

func GenerateReadme(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/api/readme" {
		err := pkg.ErrorMessage(http.StatusNotFound)
		config.RenderTemplate(w, "error", model.PageData{Error: model.ErrorData{Code: err["code"], Message: err["msg"]}})
		fmt.Println(http.StatusNotFound, err["msg"])
		return
	}

	w.Header().Set("Content-Type", "application/json")

	var body struct {
		Name        string   `json:"name"`
		Description string   `json:"description"`
		Author      string   `json:"author"`
		GitHub      string   `json:"github"`
		License     string   `json:"license"`
		Language    string   `json:"language"`
		Sections    []string `json:"sections"`
		Badges      []string `json:"badges"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.Name == "" {
		err1 := pkg.ErrorMessage(http.StatusBadRequest)
		config.RenderTemplate(w, "error", model.PageData{Error: model.ErrorData{Code: err1["code"], Message: "project name required"}})
		fmt.Println(http.StatusNotFound, "project name required")
		return
	}

	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("# %s\n\n", body.Name))
	if body.Description != "" {
		sb.WriteString(fmt.Sprintf("> %s\n\n", body.Description))
	}

	if len(body.Badges) > 0 && body.GitHub != "" {
		parts := strings.SplitN(body.GitHub, "/", 2)
		owner, repo := "", ""
		if len(parts) == 2 {
			owner, repo = parts[0], parts[1]
		} else {
			owner = body.GitHub
			repo = body.Name
		}
		for _, badge := range body.Badges {
			switch badge {
			case "license":
				sb.WriteString(fmt.Sprintf("![License](https://img.shields.io/github/license/%s/%s)\n", owner, repo))
			case "stars":
				sb.WriteString(fmt.Sprintf("![Stars](https://img.shields.io/github/stars/%s/%s)\n", owner, repo))
			case "issues":
				sb.WriteString(fmt.Sprintf("![Issues](https://img.shields.io/github/issues/%s/%s)\n", owner, repo))
			case "forks":
				sb.WriteString(fmt.Sprintf("![Forks](https://img.shields.io/github/forks/%s/%s)\n", owner, repo))
			case "go":
				sb.WriteString("![Go Version](https://img.shields.io/badge/Go-1.23-blue)\n")
			case "node":
				sb.WriteString("![Node Version](https://img.shields.io/badge/Node-20-green)\n")
			case "docker":
				sb.WriteString(fmt.Sprintf("![Docker Pulls](https://img.shields.io/docker/pulls/%s/%s)\n", owner, repo))
			}
		}
		sb.WriteString("\n")
	}

	// Sections
	for _, section := range body.Sections {
		switch section {
		case "overview":
			sb.WriteString("## Overview\n\n")
			sb.WriteString(fmt.Sprintf("%s is a ...\n\n", body.Name))

		case "features":
			sb.WriteString("## Features\n\n")
			sb.WriteString("- ✅ Feature one\n")
			sb.WriteString("- ✅ Feature two\n")
			sb.WriteString("- ✅ Feature three\n\n")

		case "installation":
			sb.WriteString("## Installation\n\n")
			sb.WriteString("```bash\n")
			if body.Language == "go" {
				sb.WriteString(fmt.Sprintf("go install github.com/%s/%s@latest\n", body.GitHub, strings.ToLower(body.Name)))
			} else if body.Language == "node" {
				sb.WriteString(fmt.Sprintf("npm install %s\n", strings.ToLower(body.Name)))
			} else {
				sb.WriteString(fmt.Sprintf("git clone https://github.com/%s/%s.git\n", body.GitHub, strings.ToLower(body.Name)))
				sb.WriteString(fmt.Sprintf("cd %s\n", strings.ToLower(body.Name)))
			}
			sb.WriteString("```\n\n")

		case "usage":
			sb.WriteString("## Usage\n\n")
			sb.WriteString("```bash\n")
			sb.WriteString(fmt.Sprintf("# Basic usage\n./%s\n", strings.ToLower(body.Name)))
			sb.WriteString("```\n\n")

		case "docker":
			sb.WriteString("## Docker\n\n")
			sb.WriteString("```bash\n")
			sb.WriteString(fmt.Sprintf("docker run -d \\\n  --name %s \\\n  -p 8080:8080 \\\n  %s/%s:latest\n", strings.ToLower(body.Name), body.GitHub, strings.ToLower(body.Name)))
			sb.WriteString("```\n\n")

		case "api":
			sb.WriteString("## API Reference\n\n")
			sb.WriteString("### Endpoint\n\n")
			sb.WriteString("```http\nGET /api/example\n```\n\n")
			sb.WriteString("| Parameter | Type | Description |\n")
			sb.WriteString("|-----------|------|-------------|\n")
			sb.WriteString("| `param`   | `string` | Description |\n\n")

		case "contributing":
			sb.WriteString("## Contributing\n\n")
			sb.WriteString("Contributions are welcome! Please follow these steps:\n\n")
			sb.WriteString("1. Fork the repository\n")
			sb.WriteString("2. Create a feature branch (`git checkout -b feat/my-feature`)\n")
			sb.WriteString("3. Commit your changes (`git commit -m 'feat: add my feature'`)\n")
			sb.WriteString("4. Push to the branch (`git push origin feat/my-feature`)\n")
			sb.WriteString("5. Open a Pull Request\n\n")

		case "license":
			sb.WriteString("## License\n\n")
			license := body.License
			if license == "" {
				license = "MIT"
			}
			sb.WriteString(fmt.Sprintf("This project is licensed under the [%s License](LICENSE).\n\n", license))

		case "author":
			if body.Author != "" {
				sb.WriteString("## Author\n\n")
				if body.GitHub != "" {
					parts := strings.SplitN(body.GitHub, "/", 2)
					sb.WriteString(fmt.Sprintf("**%s** — [@%s](https://github.com/%s)\n\n", body.Author, parts[0], parts[0]))
				} else {
					sb.WriteString(fmt.Sprintf("**%s**\n\n", body.Author))
				}
			}
		}
	}

	// Footer
	if body.Author != "" {
		sb.WriteString(fmt.Sprintf("---\n\nMade with ❤️ by **%s**\n", body.Author))
	}

	err := json.NewEncoder(w).Encode(map[string]string{"content": sb.String()})
	if err != nil {
		err1 := pkg.ErrorMessage(http.StatusInternalServerError)
		config.RenderTemplate(w, "error", model.PageData{Error: model.ErrorData{Code: err1["code"], Message: err1["msg"]}})
		fmt.Println(err1)
	}
}
