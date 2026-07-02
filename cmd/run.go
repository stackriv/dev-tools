package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/stackriv/dev-tools/internal/business/controller"
	"github.com/stackriv/dev-tools/internal/config"
)

func run(args []string) {
	if len(args) == 0 {
		var appConfig config.Config // Configuration elements

		templateCache, err := config.CreateTemplateCache() // Cache creation to switch pages
		if err != nil {
			panic(err)
		}

		appConfig.TemplateCache = templateCache // Cache view creation
		appConfig.Port = "8090"
		appConfig.StaticDir = "assets"
		config.NewAppConfig(&appConfig) // Inject config elements

		// Parsing statics files
		statics := http.FileServer(http.Dir(appConfig.StaticDir))
		http.Handle("/statics/", http.StripPrefix("/statics/", statics))

		http.HandleFunc("/", controller.Home)
		http.HandleFunc("/gitignore", controller.Gitignore)
		http.HandleFunc("/api/gitignore", controller.GenerateGitignore)
		http.HandleFunc("/license", controller.License)
		http.HandleFunc("/api/license", controller.GenerateLicense)
		http.HandleFunc("/packagejson", controller.PackageJSON)
		http.HandleFunc("/api/packagejson", controller.GeneratePackageJSON)
		http.HandleFunc("/env", controller.Env)
		http.HandleFunc("/api/env", controller.GenerateEnv)
		http.HandleFunc("/dockerfile", controller.DockerfileHandler)
		http.HandleFunc("/api/dockerfile", controller.GenerateDockerfile)
		http.HandleFunc("/uuid", controller.UUID)
		http.HandleFunc("/api/uuid", controller.GenerateUUID)
		http.HandleFunc("/base64", controller.Base64Handler)
		http.HandleFunc("/api/base64", controller.ProcessBase64)
		http.HandleFunc("/hash", controller.Hash)
		http.HandleFunc("/api/hash", controller.GenerateHash)
		http.HandleFunc("/jwt", controller.JWT)
		http.HandleFunc("/api/jwt", controller.DecodeJWT)
		http.HandleFunc("/regex", controller.Regex)
		http.HandleFunc("/api/regex", controller.TestRegex)
		http.HandleFunc("/cron", controller.Cron)
		http.HandleFunc("/api/cron", controller.ValidateCron)
		http.HandleFunc("/dns", controller.DNS)
		http.HandleFunc("/api/dns", controller.LookupDNS)
		http.HandleFunc("/ssl", controller.SSL)
		http.HandleFunc("/api/ssl", controller.CheckSSL)
		http.HandleFunc("/readme", controller.Readme)
		http.HandleFunc("/api/readme", controller.GenerateReadme)
		http.HandleFunc("/compose", controller.Compose)
		http.HandleFunc("/api/compose", controller.GenerateCompose)
		http.HandleFunc("/uptime", controller.Uptime)
		http.HandleFunc("/api/uptime", controller.CheckUptime)
		http.HandleFunc("/logs", controller.Logs)
		http.HandleFunc("/api/containers", controller.GetContainers)
		http.HandleFunc("/api/logs", controller.GetLogs)
		http.HandleFunc("/markdown", controller.Markdown)
		http.HandleFunc("/api/markdown", controller.RenderMarkdown)
		http.HandleFunc("/invoice", controller.Invoice)
		http.HandleFunc("/api/invoice", controller.GenerateInvoice)

		fmt.Println("Stackriv Dev Tools running on port http://127.0.0.1:" + appConfig.Port)
		err = http.ListenAndServe(":"+appConfig.Port, nil)
		if err != nil {
			panic(err)
			return
		}
		os.Exit(0)
	}
}
