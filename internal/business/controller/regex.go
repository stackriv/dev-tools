package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"

	"github.com/stackriv/dev-tools/internal/business/model"
	"github.com/stackriv/dev-tools/internal/config"
	"github.com/stackriv/dev-tools/internal/pkg"
)

func Regex(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		if r.URL.Path != "/regex" {
			err := pkg.ErrorMessage(http.StatusNotFound)
			config.RenderTemplate(w, "error", model.PageData{Error: model.ErrorData{Code: err["code"], Message: err["msg"]}})
			fmt.Println(http.StatusNotFound, err["msg"])
			return
		}

		config.RenderTemplate(w, "regex", model.PageData{
			Title:       "Regex Tester",
			Description: "Test your regular expressions against a string in real time.",
			Page:        "regex",
		})
	}
}

func TestRegex(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/api/regex" {
		err := pkg.ErrorMessage(http.StatusNotFound)
		config.RenderTemplate(w, "error", model.PageData{Error: model.ErrorData{Code: err["code"], Message: err["msg"]}})
		fmt.Println(http.StatusNotFound, err["msg"])
		return
	}

	w.Header().Set("Content-Type", "application/json")

	var body struct {
		Pattern string `json:"pattern"`
		Input   string `json:"input"`
		Flags   struct {
			Global      bool `json:"global"`
			Insensitive bool `json:"insensitive"`
			Multiline   bool `json:"multiline"`
		} `json:"flags"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.Pattern == "" {
		err := pkg.ErrorMessage(http.StatusBadRequest)
		config.RenderTemplate(w, "error", model.PageData{Error: model.ErrorData{Code: err["code"], Message: "pattern required"}})
		fmt.Println("pattern required")
		return
	}

	// Build the pattern with flags
	pattern := body.Pattern
	prefix := "(?"
	if body.Flags.Insensitive {
		prefix += "i"
	}
	if body.Flags.Multiline {
		prefix += "m"
	}
	if prefix != "(?" {
		pattern = prefix + ")" + pattern
	}

	re, err := regexp.Compile(pattern)
	if err != nil {
		err1 := pkg.ErrorMessage(http.StatusBadRequest)
		config.RenderTemplate(w, "error", model.PageData{Error: model.ErrorData{Code: err1["code"], Message: "invalid regex: " + err.Error()}})
		fmt.Println(http.StatusBadRequest, "invalid regex: "+err.Error())
		return
	}

	var matches []model.RegexMatch

	process := func(loc []int) {
		if loc == nil {
			return
		}
		match := body.Input[loc[0]:loc[1]]
		groups := make([]string, 0, (len(loc)-2)/2)
		for i := 2; i < len(loc); i += 2 {
			if loc[i] >= 0 {
				groups = append(groups, body.Input[loc[i]:loc[i+1]])
			} else {
				groups = append(groups, "")
			}
		}
		matches = append(matches, model.RegexMatch{
			Match:  match,
			Groups: groups,
			Start:  loc[0],
			End:    loc[1],
		})
	}

	if body.Flags.Global {
		for _, loc := range re.FindAllStringSubmatchIndex(body.Input, -1) {
			process(loc)
		}
	} else {
		if loc := re.FindStringSubmatchIndex(body.Input); loc != nil {
			process(loc)
		}
	}

	err1 := json.NewEncoder(w).Encode(map[string]interface{}{
		"valid":   true,
		"matches": matches,
		"count":   len(matches),
	})
	if err1 != nil {
		err := pkg.ErrorMessage(http.StatusInternalServerError)
		config.RenderTemplate(w, "error", model.PageData{Error: model.ErrorData{Code: err["code"], Message: err["msg"]}})
		fmt.Println(err)
	}
}
