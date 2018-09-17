package main

import (
	"io/ioutil"
	"flag"
	"fmt"
	"net/http"
	"url-shortener-golang/urlshort"
)

func main() {
	ymlFlag := flag.String("yml", "file.yml", "yml file with paths")
	jsonFlag := flag.String("json", "file.json", "json file with paths")
	flag.Parse()

	mux := defaultMux()
	pathsToUrls := map[string]string {
		"/urlshort-godoc": "https://godoc.org/github.com/gophercises/urlshort",
		"/yaml-godoc":     "https://godoc.org/gopkg.in/yaml.v2",
	}
	mapHandler := urlshort.MapHandler(pathsToUrls, mux)

	yaml, err := ioutil.ReadFile(*ymlFlag)
	if (err != nil) {
		panic(err)
	}
	yamlHandler, err := urlshort.GeneralHandler("yaml", []byte(yaml), mapHandler)
	if err != nil {
		panic(err)
	}

	json, err := ioutil.ReadFile(*jsonFlag)
	if (err != nil) {
		panic(err)
	}
	jsonHandler, err := urlshort.GeneralHandler("json", []byte(json), yamlHandler)
	if err != nil {
		panic(err)
	}

	fmt.Println("Starting the server on :7000")
	http.ListenAndServe(":7000", jsonHandler)
}

func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", hello)
	return mux
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, world!")
}