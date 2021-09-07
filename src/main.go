package main

import (
	"fmt"
	"log"
	"net/http"
)

const (
	defaultHTTPPort = 3333
	// defaultAPIRoot's the regular API root, which would get overridden via an environment variable for mocking purposes
	defaultAPIRoot = "http://interview-api.snackable.ai/api/file/"
)

func main() {
	http.Handle("/overview/", endpointHandler{createClient(defaultAPIRoot)})

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", defaultHTTPPort), nil))
}
