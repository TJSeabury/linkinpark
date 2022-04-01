package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type message struct {
	uuid    string
	message string
}

type dispatcher struct {
	jobs map[string]job
}

func (d *dispatcher) Start(rw http.ResponseWriter, r *http.Request) {
	URL := r.URL.Query().Get("url")
	if URL == "" {
		log.Println("missing URL argument")
		return
	}

	domain := getDomain(URL)
	URL = "https://" + domain

	j := NewJob(URL)
	d.jobs[j.uuid] = j
	go j.crawl()

	m := message{
		uuid:    j.uuid,
		message: "Crawling job started.",
	}
	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusCreated)
	json.NewEncoder(rw).Encode(m)
}

func (d *dispatcher) Check(rw http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var m message
	err := decoder.Decode(&m)
	check(err)
	log.Println(m.uuid)

	j, ok := d.jobs[m.uuid]
	if !ok {
		m.message = "No job with that UUID found."
		rw.Header().Set("Content-Type", "application/json")
		rw.WriteHeader(http.StatusNotFound)
		json.NewEncoder(rw).Encode(m)
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	responseJob := job{
		uuid:         j.uuid,
		status:       j.status,
		linksFound:   j.linksFound,
		linksCrawled: j.linksCrawled,
	}
	json.NewEncoder(rw).Encode(responseJob)
}

func (d *dispatcher) Finish(rw http.ResponseWriter, r *http.Request) {
	filename := domain + ".csv"

	report := creatReport(root, filename)

	log.Println("Serving", filename, " . . . ")
	rw.Header().Add("Content-Type", "text/csv")
	rw.Header().Add(
		"Content-Disposition",
		`attachment; filename="`+filename+`"`,
	)
	rw.Write(report)
}