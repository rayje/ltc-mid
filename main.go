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
		responses := requestor.MakeRequests(path)

		// Set body and Echo Headers
		content = responses[0].Body
		for k, v := range responses[0].Headers {
			w.Header().Add("X--" + k + "--X", v[0])
		}

		numResponses := len(responses)
		durations := make([]string, numResponses)
		for i := 0; i < numResponses; i++ {
			durations[i] = fmt.Sprintf("%d", responses[i].Duration.Nanoseconds())
		}
		w.Header().Add("Durations", strings.Join(durations,","))
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