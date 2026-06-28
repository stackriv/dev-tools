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

		fmt.Println("Server is running on port http://127.0.0.1:" + appConfig.Port)
		err = http.ListenAndServe(":"+appConfig.Port, nil)
		if err != nil {
			panic(err)
			return
		}
		os.Exit(0)
	}
}
