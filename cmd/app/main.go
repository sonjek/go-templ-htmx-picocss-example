package main

import (
	"github.com/sonjek/go-templ-htmx-picocss-example/internal/web"
)

func main() {
	webServer := web.NewWebServer()
	webServer.Start()
}
