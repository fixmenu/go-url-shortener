package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	mux := defaultMux()

	pathsToUrls := map[string]string{
		"/urlshort-godoc": "https://godoc.org/github.com/gophercises/urlshort",
		"/yaml-godoc":     "https://godoc.org/gopkg.in/yaml.v2",
	}

	pathHandler := mapHandler(pathsToUrls, mux)

	if err := http.ListenAndServe(":8080", pathHandler); err != nil {
		log.Fatal(err)
	}
}

func mapHandler(paths map[string]string, fallback http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path

		if dest, ok := paths[path]; ok {
			http.Redirect(w, r, dest, http.StatusFound)
			return
		}

		fallback.ServeHTTP(w, r)
	}
}

func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", greetHandler)
	return mux
}

func greetHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Greetings!")
}
