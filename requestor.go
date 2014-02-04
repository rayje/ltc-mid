package main

import (
	"net/http"
	"time"
	"fmt"
	"io/ioutil"
)

var client = &http.Client{}

type Response struct {
	Body string
	Headers http.Header
	Duration time.Duration
}

type Requestor struct {
	Config Config
	ApigeeToken string
	Apigee ApigeeConfig
}

func NewRequestor(config *Config, apigeeToken string) (Requestor, error) {
	var apigeeToken string
	var requestor Requestor
	var err error

	return Requestor{
		Config: *config,
		ApigeeToken: apigeeToken,
		Apigee: config.Apigee,
	}, nil
}

func (r *Requestor) NewRequest(url string) (http.Request, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return *req, err
	}

	if r.ApigeeToken != "" {
		req.Header.Set("Authorization", "Bearer " + r.ApigeeToken)
	}

	return *req, err
}

func (r *Requestor) BuildRequests(route string) []http.Request {
	requests := make([]http.Request, len(r.Config.EndPoints))
	var err error

	for i := 0; i < cap(requests); i++ {
		url := r.Config.EndPoints[i].Url + route
		requests[i], err = r.NewRequest(url)
		if err != nil {
			fmt.Println("Error building requests")
			fmt.Println(err)
		}
	}

	return requests
}

func (r *Requestor) MakeRequests(route string) []Response {
	numRequests := len(r.Config.EndPoints)
	requests := r.BuildRequests(route)
	res := make(chan Response, numRequests)
	results := make([]Response, numRequests)

	for i := 0; i < numRequests; i++ {
		go runRequest(&requests[i], res)
	}

	for i := 0; i < cap(res); i++ {
		results[i] = <-res
	}
	close(res)

	return results
}

func runRequest(req *http.Request, res chan Response) {
	start := time.Now()
	r, err := client.Do(req)

	response := Response{
		Duration:  time.Since(start),
	}

	if err != nil {
		fmt.Println(err)
	} else {
		if body, err := ioutil.ReadAll(r.Body); err != nil {
			fmt.Println(err)
		} else {
			if r.StatusCode < 200 || r.StatusCode >= 300 {
				fmt.Println("======================================")
				fmt.Println("Status: " + r.Status)
				for k, v := range r.Header {
					fmt.Println(k, ":", v)
				}
				fmt.Println(string(body))
				fmt.Println("======================================")
			} else {
				response.Body = string(body)
				response.Headers = r.Header
			}
		}
	}

	res <- response
}