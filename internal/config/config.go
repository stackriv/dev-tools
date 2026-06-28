package config

import (
	"text/template"
)

type Config struct {
	TemplateCache map[string]*template.Template // Cache
	Port          string
	StaticDir     string
}
