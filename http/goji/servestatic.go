package goji

import (
	"fmt"
	"log"
	"net/http"

	"github.com/zenazn/goji"
)

// ServeStaticFiles asks Goji to serve static files in the specified directory
func ServeStaticFiles(prefix, directorypath string) {
	log.Printf("Serving static files in %s", directorypath)
	goji.Handle(fmt.Sprintf("%s*", prefix), http.StripPrefix(prefix, http.FileServer(http.Dir(directorypath))))
}
