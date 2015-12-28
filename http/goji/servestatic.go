package goji

import (
	"fmt"
	"log"
	"net/http"

	"github.com/zenazn/goji/web"
)

// ServeStaticFiles asks Goji to serve static files in the specified directory
func ServeStaticFiles(prefix, directorypath string, mux *web.Mux) {
	log.Printf("Serving static files in %s", directorypath)
	mux.Handle(fmt.Sprintf("%s*", prefix), http.StripPrefix(prefix, http.FileServer(http.Dir(directorypath))))
}
