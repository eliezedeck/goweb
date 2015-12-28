package validation

import (
	"fmt"
	"log"
	"strings"

	"github.com/xeipuuv/gojsonschema"
)

var schemaCache map[string]*gojsonschema.Schema

func init() {
	schemaCache = make(map[string]*gojsonschema.Schema)
}

// LoadSchema loads a JSON Schema file from `path` and assigns it to the `schemaid`.
// Note that any error here is fatal.
func LoadSchema(path, schemaid string) {
	var loader gojsonschema.JSONLoader

	loader = gojsonschema.NewReferenceLoader(path)
	schema, err := gojsonschema.NewSchema(loader)
	if err != nil {
		log.Fatalf("Failed to load Schema for '%s': %s", schemaid, err)
	}

	schemaCache[schemaid] = schema
	log.Printf("Loaded Schema '%s'", schemaid)
}

// Validate validates `rawdata` using the cached `schemaid` and returns a
// user-friendly error message (as a string) in case of any validation error
func Validate(schemaid string, rawdata interface{}) (bool, string) {
	schema, ok := schemaCache[schemaid]
	if !ok {
		return false, "Schema not found"
	}

	jloader := gojsonschema.NewGoLoader(rawdata)
	result, err := schema.Validate(jloader)
	if err != nil {
		return false, "Validation error"
	}

	if result.Valid() {
		return true, ""
	}

	firsterror := result.Errors()[0]
	return false, formatValidationError(firsterror)
}

func formatValidationError(e gojsonschema.ResultError) string {
	return fmt.Sprintf("%s: %s", strings.Title(e.Field()), e.Description())
}
