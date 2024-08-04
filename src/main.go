package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"main/src/urlshort"
	"net/http"
	"os"
)

func main() {
	mux := defaultMux()

	// Build the MapHandler using the mux as the fallback
	pathsToUrls := map[string]string{
		"/urlshort-godoc": "https://godoc.org/github.com/gophercises/urlshort",
		"/yaml-godoc":     "https://godoc.org/gopkg.in/yaml.v2",
	}
	mapHandler := urlshort.MapHandler(pathsToUrls, mux)

	// Build the YAMLHandler using the mapHandler as the
	// fallback
	yamlFileLocation := flag.String("yaml", "whatever.yaml", "Yaml with redirections")
	flag.Parse()

	yamlHandler, err := urlshort.YAMLHandler(getFileBytes(*yamlFileLocation), mapHandler)
	if err != nil {
		panic(err)
	}
	fmt.Println("Starting the server on :8080")
	http.ListenAndServe(":8080", yamlHandler)
}

func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", hello)
	return mux
}

func hello(w http.ResponseWriter, r *http.Request) {
}

func getFileBytes(fileName string) []byte { //https://gobyexample.com/reading-files
	f, err := os.Open(fileName)
	if err != nil {
		log.Fatalf("Could not open file %s", fileName)
	}

	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(f)
	if err != nil {
		log.Fatalf("Could not read file %s", fileName)
	}

	return buf.Bytes()
}
