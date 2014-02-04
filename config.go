package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
)

type ApigeeConfig struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Apikey   string `json:"apikey"`
	Apiurl   string `json:"apiurl"`
}

type EndPoint struct {
	Url  string `json:"url"`
}

type EndPoints []EndPoint

type Config struct {
	EndPoints       EndPoints      `json:"endpoints"`
	Port			string         `json:"port"`
	Apigee          ApigeeConfig   `json:"apigee"`
}

func getConfig() Config {
	configFile := flag.String("config", "config.json", "Location of config file")
	flag.Parse()

	file, err := ioutil.ReadFile(*configFile)
	if err != nil {
		fmt.Printf("File error: %v\n", err)
		os.Exit(1)
	}

	var config Config
	err = json.Unmarshal(file, &config)
	if err != nil {
		fmt.Printf("JSON error: %v\n", err)
		os.Exit(1)
	}

	return config
}