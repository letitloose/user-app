package static

import (
	"net/http"
)

func AddHandlersToMux(mux *http.ServeMux) {
	fileServer := http.FileServer(http.Dir("./pkg/static/"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))
}
