package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func decode(r *http.Request) message {
	decoder := json.NewDecoder(r.Body)
	var m message
	err := decoder.Decode(&m)
	check(err)
	return m
}

type message struct {
	uuid    string
	message string
}

type dispatcher struct {
	jobs map[string]job
}

func (d *dispatcher) Start(rw http.ResponseWriter, r *http.Request) {
	m := decode(r)

	// message is the domain to crawl in this case.
	// Definitelly a code smell that this requires a comment.
	URL := m.message
	if URL == "" || !IsUrl(URL) {
		log.Println("Bad URL argument!")
		return
	}

	domain := getDomain(URL)
	URL = "https://" + domain

	j := NewJob(URL)
	d.jobs[j.uuid] = j
	log.Println("Started", j.uuid)
	log.Println("Crawling", j.domain, " . . . ")
	go j.crawl()

	m = message{
		uuid:    j.uuid,
		message: "Crawling job started.",
	}
	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusCreated)
	json.NewEncoder(rw).Encode(m)
}

func (d *dispatcher) Check(rw http.ResponseWriter, r *http.Request) {
	m := decode(r)
	log.Println("Check", m.uuid)

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
	m := decode(r)
	log.Println("Finish", m.uuid)

	report := creatReport(root, filename)

	log.Println("Serving", filename, " . . . ")
	rw.Header().Add("Content-Type", "text/csv")
	rw.Header().Add(
		"Content-Disposition",
		`attachment; filename="`+filename+`"`,
	)
	rw.Write(report)
}
