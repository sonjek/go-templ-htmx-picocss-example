package handlers

import (
	"fmt"
	"net/http"

	"github.com/a-h/templ"
	"github.com/sonjek/go-templ-htmx-picocss-example/internal/web/templ/components"
	"github.com/sonjek/go-templ-htmx-picocss-example/internal/web/templ/page"
	"github.com/sonjek/go-templ-htmx-picocss-example/internal/web/templ/view"
)

type Handlers struct{}

func NewHandler() *Handlers {
	return &Handlers{}
}

func handleRenderError(err error) {
	if err != nil {
		fmt.Println("Render error: ", err)
	}
}

// Set header and render error message
func sendErrorMsg(w http.ResponseWriter, r *http.Request, errorMsg string) {
	w.WriteHeader(http.StatusBadRequest)
	handleRenderError(components.ErrorMsg(errorMsg).Render(r.Context(), w))
}

func (h *Handlers) Page404(w http.ResponseWriter, r *http.Request) {
	templ.Handler(
		page.Index(view.NotFoundComponent()),
		templ.WithStatus(http.StatusNotFound),
	).ServeHTTP(w, r)
}
