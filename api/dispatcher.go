package main

import (
	"encoding/json"
	"net/http"
	"strings"
)

type message struct {
	Uuid    string `json:"uuid"`
	Message string `json:"message"`
}

type dispatcher struct {
	jobs map[string]*job
}

func NewDispatcher() dispatcher {
	return dispatcher{
		jobs: make(map[string]*job),
	}
}

func (d *dispatcher) Start(rw http.ResponseWriter, r *http.Request) {
	m := decode(r)

	// message is the domain to crawl in this case.
	// Definitelly a code smell that this requires a comment.
	URL := m.Message

	domain, err := getDomain(URL)

	if err != nil || URL == "" || !IsUrl(URL) {
		m.Message = "Bad URL argument! " + URL + " > "
		//log.Println("Bad URL argument!")
		rw.Header().Set("Content-Type", "application/json")
		rw.WriteHeader(http.StatusUnprocessableEntity)
		json.NewEncoder(rw).Encode(m)
		return
	}

	j := NewJob(domain)
	d.jobs[j.Uuid] = &j

	m = message{
		Uuid:    j.Uuid,
		Message: "Crawling job started.",
	}
	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusCreated)
	json.NewEncoder(rw).Encode(m)

	go j.Start()

}

func (d *dispatcher) Check(rw http.ResponseWriter, r *http.Request) {
	m := decode(r)

	j, ok := d.jobs[m.Uuid]
	if !ok {
		m.Message = "No job with that UUID found."
		rw.Header().Set("Content-Type", "application/json")
		rw.WriteHeader(http.StatusNotFound)
		json.NewEncoder(rw).Encode(m)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	responseJob := job{
		Uuid:         j.Uuid,
		Status:       j.Status,
		LinksFound:   j.LinksFound,
		LinksCrawled: j.LinksCrawled,
	}
	json.NewEncoder(rw).Encode(responseJob)
}

func (d *dispatcher) Finish(rw http.ResponseWriter, r *http.Request) {
	var m message
	m.Uuid = strings.TrimPrefix(r.URL.Path, "/api/finish/")
	//log.Println("Getting job:", m.Uuid)

	j, ok := d.jobs[m.Uuid]
	if !ok {
		m.Message = "No job with that UUID found."
		rw.Header().Set("Content-Type", "application/json")
		rw.WriteHeader(http.StatusNotFound)
		json.NewEncoder(rw).Encode(m)
		return
	}

	if j.Status != "done" {
		m.Message = "That job is not yet done."
		rw.Header().Set("Content-Type", "application/json")
		rw.WriteHeader(http.StatusAccepted)
		json.NewEncoder(rw).Encode(m)
		return
	}

	filename := j.Domain + ".csv"

	report := creatReport(j.Data, filename)

	//log.Println("Serving", filename, " . . . ")
	rw.Header().Add("Content-Type", "text/csv")
	rw.Header().Add(
		"Content-Disposition",
		`attachment; filename="`+filename+`"`,
	)
	rw.Write(report)

	// delete data to save memory
	delete(d.jobs, m.Uuid)
}
