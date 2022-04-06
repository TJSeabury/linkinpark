package main

import (
	"log"

	"github.com/gocolly/colly"
	"github.com/google/uuid"
)

type job struct {
	Uuid          string
	Status        string
	StatusChannel chan string
	Domain        string
	LinksFound    int
	LinksCrawled  int
	Data          map[string]pageInfo
	DataChannel   chan map[string]pageInfo
}

func NewJob(domain string) job {
	id := uuid.New()
	return job{
		Uuid:          "job_" + id.String(),
		Status:        "starting",
		StatusChannel: make(chan string),
		Domain:        domain,
		DataChannel:   make(chan map[string]pageInfo),
	}
}

func (j *job) addLinksFound(n int) {
	j.LinksFound += n
}

func (j *job) addLinksCrawled(n int) {
	j.LinksCrawled += n
}

func (j *job) crawl() {
	// Instantiate default collector
	c := colly.NewCollector(
		// MaxDepth is 2, so only the links on the scraped page
		// and links on those pages are visited
		colly.MaxDepth(1),
		colly.AllowedDomains(j.Domain, "www."+j.Domain),
		colly.Async(true),
	)

	// Limit the maximum parallelism to 2
	// This is necessary if the goroutines are dynamically
	// created to control the limit of simultaneous requests.
	//
	// Parallelism can be controlled also by spawning fixed
	// number of go routines.
	c.Limit(&colly.LimitRule{DomainGlob: "*", Parallelism: 8})

	log.Println("Started", j.Uuid)
	log.Println("Crawling", j.Domain, " . . . ")

	go crawl(
		j,
		c,
		j.Domain,
		make(map[string]pageInfo),
	)
}

func (j *job) checkData() *job {
	j.Status = <-j.StatusChannel
	j.Data = <-j.DataChannel
	return j
}
