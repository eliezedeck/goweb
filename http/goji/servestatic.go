package goji

import (
	"fmt"
	"log"
	"net/http"

	"goji.io"
	"goji.io/pat"
)

// ServeStaticFiles asks Goji to serve static files in the specified directory
func ServeStaticFiles(prefix, directorypath string, mux *goji.Mux) {
	log.Printf("Serving static files in %s", directorypath)
	mux.Handle(pat.Get(fmt.Sprintf("%s*", prefix)), http.StripPrefix(prefix, http.FileServer(http.Dir(directorypath))))
}
