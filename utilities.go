package main

import (
	"encoding/json"
	"errors"
	"log"
	"net"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func IsUrl(str string) bool {
	u, err := url.ParseRequestURI(str)
	return err == nil && u.Scheme != "" && u.Host != ""
}

func IsHostReachable(host string) bool {
	timeout := 1 * time.Second
	conn, err := net.DialTimeout("tcp", host+":80", timeout)
	if err != nil {
		//conn.Close()
		log.Println(err)
		return false
	}
	conn.Close()
	return true
}

func IsDomainValid(url string) bool {
	reg := regexp.MustCompile("^(((?!\\-))(xn\\-\\-)?[a-z0-9\\-_]{0,61}[a-z0-9]{1,1}\\.)*(xn\\-\\-)?([a-z0-9\\-]{1,61}|[a-z0-9\\-]{1,30})\\.[a-z]{2,}$")
	return reg.Match([]byte(url))
}

func getDomain(uri string) (string, error) {
	u, err := url.Parse(uri)
	if err != nil {
		return "", err
	}
	parts := strings.Split(u.Hostname(), ".")
	if len(parts)-2 < 0 || len(parts)-1 < 0 {
		return "", errors.New("Invalid domain!")
	}
	return parts[len(parts)-2] + "." + parts[len(parts)-1], nil
}

func decode(r *http.Request) message {
	decoder := json.NewDecoder(r.Body)
	var m message
	err := decoder.Decode(&m)
	check(err)
	return m
}
