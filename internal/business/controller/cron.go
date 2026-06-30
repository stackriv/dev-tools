package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/stackriv/dev-tools/internal/business/model"
	"github.com/stackriv/dev-tools/internal/config"
	"github.com/stackriv/dev-tools/internal/pkg"
)

func Cron(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		if r.URL.Path != "/cron" {
			err := pkg.ErrorMessage(http.StatusNotFound)
			config.RenderTemplate(w, "error", model.PageData{Error: model.ErrorData{Code: err["code"], Message: err["msg"]}})
			fmt.Println(http.StatusNotFound, err["msg"])
			return
		}

		config.RenderTemplate(w, "cron", model.PageData{
			Title:       "Cron Expression Builder",
			Description: "Build and validate cron expressions visually.",
			Page:        "cron",
		})
	}
}

func ValidateCron(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/api/cron" {
		err := pkg.ErrorMessage(http.StatusNotFound)
		config.RenderTemplate(w, "error", model.PageData{Error: model.ErrorData{Code: err["code"], Message: err["msg"]}})
		fmt.Println(http.StatusNotFound, err["msg"])
		return
	}

	w.Header().Set("Content-Type", "application/json")

	var body struct {
		Expression string `json:"expression"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.Expression == "" {
		err := pkg.ErrorMessage(http.StatusBadRequest)
		config.RenderTemplate(w, "error", model.PageData{Error: model.ErrorData{Code: err["code"], Message: "expression required"}})
		fmt.Println(err)
	}

	parts := strings.Fields(body.Expression)
	if len(parts) != 5 {
		err := pkg.ErrorMessage(http.StatusBadRequest)
		config.RenderTemplate(w, "error", model.PageData{Error: model.ErrorData{Code: err["code"], Message: "cron expression must have 5 fields: minute hour day month weekday"}})
		fmt.Println(err)
	}

	description := describeCron(parts)
	nextRuns := nextCronRuns(parts, 5)

	err := json.NewEncoder(w).Encode(map[string]interface{}{
		"valid":       true,
		"expression":  body.Expression,
		"description": description,
		"next_runs":   nextRuns,
	})
	if err != nil {
		err := pkg.ErrorMessage(http.StatusInternalServerError)
		config.RenderTemplate(w, "error", model.PageData{Error: model.ErrorData{Code: err["code"], Message: err["msg"]}})
		fmt.Println(err)
	}
}

func describeCron(parts []string) string {
	minute, hour, day, month, weekday := parts[0], parts[1], parts[2], parts[3], parts[4]

	months := map[string]string{
		"1": "January", "2": "February", "3": "March", "4": "April",
		"5": "May", "6": "June", "7": "July", "8": "August",
		"9": "September", "10": "October", "11": "November", "12": "December",
	}
	days := map[string]string{
		"0": "Sunday", "1": "Monday", "2": "Tuesday", "3": "Wednesday",
		"4": "Thursday", "5": "Friday", "6": "Saturday", "7": "Sunday",
	}

	var desc strings.Builder

	// Frequency
	if minute == "*" && hour == "*" && day == "*" && month == "*" && weekday == "*" {
		return "Every minute"
	}
	if minute == "0" && hour == "*" {
		desc.WriteString("Every hour")
	} else if minute != "*" && hour == "*" {
		desc.WriteString("Every hour at minute " + minute)
	} else if minute == "0" && hour != "*" {
		desc.WriteString("At " + hour + ":00")
	} else if strings.HasPrefix(minute, "*/") {
		desc.WriteString("Every " + minute[2:] + " minutes")
	} else if strings.HasPrefix(hour, "*/") {
		desc.WriteString("Every " + hour[2:] + " hours")
	} else {
		desc.WriteString("At " + hour + ":" + padZero(minute))
	}

	// Day of the week
	if weekday != "*" {
		if name, ok := days[weekday]; ok {
			desc.WriteString(" on " + name)
		} else {
			desc.WriteString(" on weekday " + weekday)
		}
	}

	// Day of the month
	if day != "*" {
		desc.WriteString(" on day " + day)
	}

	// Month
	if month != "*" {
		if name, ok := months[month]; ok {
			desc.WriteString(" in " + name)
		} else {
			desc.WriteString(" in month " + month)
		}
	}

	return desc.String()
}

func padZero(s string) string {
	if len(s) == 1 {
		return "0" + s
	}
	return s
}

func nextCronRuns(parts []string, count int) []string {
	var results []string
	now := time.Now().UTC().Truncate(time.Minute).Add(time.Minute)

	minute, hour, day, month, weekday := parts[0], parts[1], parts[2], parts[3], parts[4]

	for len(results) < count {
		if matchField(minute, now.Minute(), 0, 59) &&
			matchField(hour, now.Hour(), 0, 23) &&
			matchField(day, now.Day(), 1, 31) &&
			matchField(month, int(now.Month()), 1, 12) &&
			matchField(weekday, int(now.Weekday()), 0, 6) {
			results = append(results, now.Format("2006-01-02 15:04 UTC"))
		}
		now = now.Add(time.Minute)
		if len(results) == 0 && now.Sub(time.Now()) > 366*24*time.Hour {
			break
		}
	}

	return results
}

func matchField(field string, value, min, max int) bool {
	if field == "*" {
		return true
	}
	if strings.HasPrefix(field, "*/") {
		step, err := strconv.Atoi(field[2:])
		if err != nil {
			return false
		}
		return (value-min)%step == 0
	}
	if strings.Contains(field, ",") {
		for _, part := range strings.Split(field, ",") {
			v, err := strconv.Atoi(strings.TrimSpace(part))
			if err == nil && v == value {
				return true
			}
		}
		return false
	}
	if strings.Contains(field, "-") {
		rangeParts := strings.SplitN(field, "-", 2)
		lo, err1 := strconv.Atoi(rangeParts[0])
		hi, err2 := strconv.Atoi(rangeParts[1])
		if err1 != nil || err2 != nil {
			return false
		}
		return value >= lo && value <= hi
	}
	v, err := strconv.Atoi(field)
	return err == nil && v == value
}
