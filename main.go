package main

import (
	"fmt"
	"net/http"
	"time"
	"log"
	"strings"
)

type RequestHandler struct {
	Config Config
}

func (rh *RequestHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()
	path := r.URL.Path
	var content string

	if path != "/status" {
		token := r.Header.Get("Authorization")

		var apigeeToken string
		if (token != "") {
			apigeeToken = strings.Fields(token)[1]
		}

		requestor := NewRequestor(&rh.Config, apigeeToken)
		response := requestor.MakeRequests(path)

		// Echo Headers
		for k, v := range response[0].Headers {
			w.Header().Add("X--" + k + "--X", v[0])
		}

		content = response[0].Body
		w.Header().Add("ReadTime", time.Now().Sub(startTime).String())
	}

	fmt.Fprint(w, content)
}

func newRequestHandler(config *Config) *RequestHandler {
	return &RequestHandler{Config: *config}
}

func main() {
	config := getConfig()
	port := config.Port

	http.Handle("/", newRequestHandler(&config))
	fmt.Println("Running on port " + port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}