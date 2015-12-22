package goji

import (
	"encoding/json"
	"net/http"

	"github.com/zenazn/goji/web"
)

const maxJSONSize = 64 * 1024

// ParseJSONFormMiddleware adds a "json" entry into c.Env when there is a JSON
// format that has been detected from a POST form.
func ParseJSONFormMiddleware(c *web.C, h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" && r.ContentLength <= maxJSONSize && r.Header.Get("content-type") == "application/json" {
			decoder := json.NewDecoder(r.Body)
			var j map[string]interface{}
			if err := decoder.Decode(&j); err == nil {
				c.Env["json"] = j
			}
		}

		h.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}
