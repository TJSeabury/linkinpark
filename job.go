package main

import (
	"log"

	"github.com/gocolly/colly/v2"
	"github.com/google/uuid"
)

type job struct {
	Uuid         string
	Status       string
	Domain       string
	LinksFound   int
	LinksCrawled int
	Data         map[string]pageInfo
}

func NewJob(domain string) job {
	id := uuid.New()
	return job{
		Uuid:   "job_" + id.String(),
		Status: "starting",
		Domain: domain,
	}
}

func (j *job) addLinksFound(n int) {
	j.LinksFound += n
}

func (j *job) addLinksCrawled(n int) {
	j.LinksCrawled += n
}

func (j *job) crawl() {
	log.Println("Started", j.Uuid, "with domain", j.Domain, ". . . ")

	j.Status = "crawling"

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

	j.Data = crawl(
		j,
		c,
		j.Domain,
		make(map[string]pageInfo),
	)

}
