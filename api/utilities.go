package main

import (
	"encoding/json"
	"net"
	"net/http"
	"net/url"
	"regexp"
	"time"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func IsUrl(uri string) bool {
	u, err := url.ParseRequestURI(uri)
	return err == nil && u.Scheme != "" && u.Host != ""
}

func IsExternal(domain string, uri string) bool {
	u, err := url.ParseRequestURI(uri)
	check(err)
	return domain != u.Host
}

func IsHostReachable(host string) bool {
	timeout := 1 * time.Second
	conn, err := net.DialTimeout("tcp", host+":80", timeout)
	if err != nil {
		//conn.Close()
		//log.Println(err)
		return false
	}
	conn.Close()
	return true
}

func IsDomainValid(uri string) bool {
	reg := regexp.MustCompile("^(((?!\\-))(xn\\-\\-)?[a-z0-9\\-_]{0,61}[a-z0-9]{1,1}\\.)*(xn\\-\\-)?([a-z0-9\\-]{1,61}|[a-z0-9\\-]{1,30})\\.[a-z]{2,}$")
	return reg.Match([]byte(uri))
}

func getDomain(uri string) (string, error) {
	u, err := url.ParseRequestURI(uri)
	if err != nil {
		return "", err
	}
	return u.Host, nil
}

func decode(r *http.Request) message {
	decoder := json.NewDecoder(r.Body)
	var m message
	err := decoder.Decode(&m)
	check(err)
	return m
}

type SecurityHeader struct {
	Name  string
	Value string
}

type SecurityHeaders struct {
	StrictTransportSecurity SecurityHeader
	ContentSecurityPolicy   SecurityHeader
	XFrameOptions           SecurityHeader
	XContentTypeOptions     SecurityHeader
	ReferrerPolicy          SecurityHeader
	PermissionsPolicy       SecurityHeader
}

var SecurityHeadersEnum = SecurityHeaders{
	StrictTransportSecurity: SecurityHeader{
		Name:  "Strict-Transport-Security",
		Value: "",
	},
	ContentSecurityPolicy: SecurityHeader{
		Name:  "Content-Security-Policy",
		Value: "",
	},
	XFrameOptions: SecurityHeader{
		Name:  "X-Frame-Options",
		Value: "",
	},
	XContentTypeOptions: SecurityHeader{
		Name:  "X-Content-Type-Options",
		Value: "",
	},
	ReferrerPolicy: SecurityHeader{
		Name:  "Referrer-Policy",
		Value: "",
	},
	PermissionsPolicy: SecurityHeader{
		Name:  "Permissions-Policy",
		Value: "",
	},
}

func HasSecurityHeaders() {

}
