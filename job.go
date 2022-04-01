package main

import (
	"github.com/gocolly/colly"
	"github.com/google/uuid"
)

type job struct {
	uuid         string
	status       string
	domain       string
	linksFound   int
	linksCrawled int
	data         map[string]pageInfo
}

func NewJob(domain string) job {
	id := uuid.New()
	return job{
		uuid:   "job_" + id.String(),
		status: "starting",
		domain: domain,
	}
}

func (j *job) updateStatus(status string) {
	j.status = status
}

func (j *job) addLinksFound(n int) {
	j.linksFound += n
}

func (j *job) addLinksCrawled(n int) {
	j.linksCrawled += n
}

func (j *job) crawl() {
	// Instantiate default collector
	c := colly.NewCollector(
		// MaxDepth is 2, so only the links on the scraped page
		// and links on those pages are visited
		colly.MaxDepth(1),
		colly.AllowedDomains(j.domain, "www."+j.domain),
		colly.Async(true),
	)

	// Limit the maximum parallelism to 2
	// This is necessary if the goroutines are dynamically
	// created to control the limit of simultaneous requests.
	//
	// Parallelism can be controlled also by spawning fixed
	// number of go routines.
	c.Limit(&colly.LimitRule{DomainGlob: "*", Parallelism: 8})

	j.data = crawl(
		j,
		c,
		j.domain,
		make(map[string]pageInfo),
	)
}
