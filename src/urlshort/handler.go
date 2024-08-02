package urlshort

import (
	"net/http"

	"github.com/go-yaml/yaml"
)

// MapHandler will return an http.HandlerFunc (which also
// implements http.Handler) that will attempt to map any
// paths (keys in the map) to their corresponding URL (values
// that each key in the map points to, in string format).
// If the path is not provided in the map, then the fallback
// http.Handler will be called instead.
func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) { //http.HandlerFunc is a type that is a func :), so no need to parse
		if dest, ok := pathsToUrls[r.URL.Path]; ok { //If you find a request whose URL path matches to a key in the map, ok will be true, otherwise, false. If ok finds something, dest will have a value
			http.Redirect(w, r, dest, http.StatusFound) //Success, redirect
			return
		}
		fallback.ServeHTTP(w, r)
	}
}

// YAMLHandler will parse the provided YAML and then return
// an http.HandlerFunc (which also implements http.Handler)
// that will attempt to map any paths to their corresponding
// URL. If the path is not provided in the YAML, then the
// fallback http.Handler will be called instead.
//
// YAML is expected to be in the format:
//
//   - path: /some-path
//     url: https://www.some-url.com/demo
//
// The only errors that can be returned all related to having
// invalid YAML data.
//
// See MapHandler to create a similar http.HandlerFunc via
// a mapping of paths to urls.
func YAMLHandler(yml []byte, fallback http.Handler) (http.HandlerFunc, error) {
	//Turn YAML into Map, and then call MapHandler(), see https://pkg.go.dev/gopkg.in/yaml.v2#example-Unmarshal-Embedded

	var urlInfos []urlInfo

	err := yaml.Unmarshal(yml, &urlInfos) //parse the yaml
	if err != nil {
		return nil, err
	}

	mapUI := make(map[string]string) //create the map with the parse yaml

	for _, ui := range urlInfos {
		mapUI[ui.path] = ui.url
	}

	return MapHandler(mapUI, fallback), nil //there can be no errors as they are returned either by err or by MapHandler()
}

type urlInfo struct {
	path string `yaml:"path"`
	url  string `yaml:"url"`
}
