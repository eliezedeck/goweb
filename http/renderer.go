package http

import (
	"log"

	"github.com/unrolled/render"
)

// H is the HTML Renderer for normal pages (non-admin)
var H *render.Render

// SetupHTTPRenderer prepares the HTML/JSON/... template renderers ready for use
// before any actual rendering can take place
func SetupHTTPRenderer(viewspath, layout string) {
	log.Println("Preparing HTML renderer...")
	H = render.New(render.Options{
		Directory: viewspath,
		Layout:    layout,
		Delims: render.Delims{
			Left:  "[[",
			Right: "]]",
		},
		IndentJSON: false,
	})
}
