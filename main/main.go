package main

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"log"
	"net/http"
)

func main() {
	mux := defaultMux()

	pathsToUrls := map[string]string{
		"/urlshort-godoc": "https://godoc.org/github.com/gophercises/urlshort",
		"/yaml-godoc":     "https://godoc.org/gopkg.in/yaml.v2",
	}

	yaml := `
- path: /urlshort
  url: https://github.com/gophercises/urlshort
- path: /urlshort-final
  url: https://github.com/gophercises/urlshort/tree/solution
`

	pathHandler := MapHandler(pathsToUrls, mux)
	yamlHandler, err := YAMLHandler([]byte(yaml), pathHandler)

	if err != nil {
		panic(err)
	}

	if err := http.ListenAndServe(":8080", yamlHandler); err != nil {
		log.Fatal(err)
	}
}

func MapHandler(paths map[string]string, fallback http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path

		if dest, ok := paths[path]; ok {
			http.Redirect(w, r, dest, http.StatusFound)
			return
		}

		fallback.ServeHTTP(w, r)
	}
}

func YAMLHandler(yaml []byte, fallback http.Handler) (http.HandlerFunc, error) {
	parsedYaml, err := parseYAML()
	if err != nil {
		return nil, err
	}
	pathMap := buildMap(parsedYaml)

	return MapHandler(pathMap, fallback), nil
}

func parseYAML() ([]pathUrl, error) {
	var pathUrls []pathUrl

	buf, err := ioutil.ReadFile("./shortpaths.yaml")

	if err != nil {
		panic(err)
	}

	err = yaml.Unmarshal(buf, &pathUrls)

	if err != nil {
		return nil, err
	}

	return pathUrls, nil
}

func buildMap(pathUrls []pathUrl) map[string]string {
	pathMap := make(map[string]string)
	for _, val := range pathUrls {
		pathMap[val.Path] = val.Url
	}

	return pathMap
}

func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", greetHandler)
	return mux
}

func greetHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Greetings!")
}

type pathUrl struct {
	Path string
	Url  string
}
