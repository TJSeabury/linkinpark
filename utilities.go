package main

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"strings"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func IsUrl(str string) bool {
	u, err := url.Parse(str)
	return err == nil && u.Scheme != "" && u.Host != ""
}

func getDomain(uri string) string {
	u, err := url.Parse(uri)
	if err != nil {
		log.Fatal(err)
	}
	parts := strings.Split(u.Hostname(), ".")
	return parts[len(parts)-2] + "." + parts[len(parts)-1]
}

func decode(r *http.Request) message {
	decoder := json.NewDecoder(r.Body)
	var m message
	err := decoder.Decode(&m)
	check(err)
	return m
}
