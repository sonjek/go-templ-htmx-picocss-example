package main

import (
	"github.com/sonjek/go-templ-htmx-picocss-example/internal/web"
	"github.com/sonjek/go-templ-htmx-picocss-example/internal/web/handlers"
)

func main() {
	hs := handlers.NewHandler()
	webServer := web.NewServer(hs)
	webServer.Start()
}
