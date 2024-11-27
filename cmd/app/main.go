package main

import (
	"github.com/sonjek/go-templ-htmx-picocss-example/internal/service"
	"github.com/sonjek/go-templ-htmx-picocss-example/internal/storage"
	"github.com/sonjek/go-templ-htmx-picocss-example/internal/web"
	"github.com/sonjek/go-templ-htmx-picocss-example/internal/web/handlers"
)

func main() {
	db := storage.NewDbStorage()
	storage.DBMigrate(db)
	storage.SeedData(db)

	noteService := service.NewNoteService(db)
	handlers := handlers.NewHandler(db, noteService)
	webServer := web.NewServer(handlers)
	webServer.Start()
}
