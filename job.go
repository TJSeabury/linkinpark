package main

import (
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
	Crawler      *colly.Collector
}

func NewJob(domain string) job {
	id := uuid.New()

	// Instantiate default collector
	crawler := colly.NewCollector(
		// MaxDepth is 2, so only the links on the scraped page
		// and links on those pages are visited
		colly.MaxDepth(1),
		colly.AllowedDomains(domain, "www."+domain),
		colly.Async(),
	)

	// Limit the maximum parallelism to 2
	// This is necessary if the goroutines are dynamically
	// created to control the limit of simultaneous requests.
	//
	// Parallelism can be controlled also by spawning fixed
	// number of go routines.
	crawler.Limit(&colly.LimitRule{Parallelism: 2})

	return job{
		Uuid:    "job_" + id.String(),
		Status:  "starting",
		Domain:  domain,
		Crawler: crawler,
	}
}

func (j *job) addLinksFound(n int) {
	j.LinksFound += n
}

func (j *job) addLinksCrawled(n int) {
	j.LinksCrawled += n
}

func (j *job) Start() {
	//log.Println("Started", j.Uuid, "with domain", j.Domain, ". . . ")

	j.Status = "crawling"

	j.Data = j.crawl(
		"http://"+j.Domain,
		make(map[string]pageInfo),
	)

	j.Status = "done"

}

func (j *job) crawl(url string, pi map[string]pageInfo) map[string]pageInfo {
	//log.Println("Checking " + url)

	_, exists := pi[url]
	visited, _ := j.Crawler.HasVisited(url)
	if exists || !IsUrl(url) || visited {
		return make(map[string]pageInfo)
	}

	//log.Println("Crawling " + url)

	j.addLinksFound(1)
	j.addLinksCrawled(1)

	p := pageInfo{
		Url: url,
	}

	links := make(map[string]bool)

	j.Crawler.OnResponse(func(r *colly.Response) {
		p.StatusCode = r.StatusCode
		headers := *r.Headers
		p.ContentType = headers.Get("Content-Type")
		p.Size = len(r.Body)
	})

	j.Crawler.OnError(func(r *colly.Response, _ error) {
		//log.Println("error:", r.StatusCode, err)
		p.StatusCode = r.StatusCode
	})

	j.Crawler.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Request.AbsoluteURL(e.Attr("href"))
		if IsUrl(link) {
			p.Links++
			if _, exists := pi[link]; !exists {
				links[link] = true
			}
		}
	})

	j.Crawler.Visit(p.Url)
	j.Crawler.Wait()

	//log.Println("number of found links:", len(links))

	j.addLinksFound(len(links))

	pi[url] = p

	for link := range links {
		res := j.crawl(link, pi)
		for k, v := range res {
			pi[k] = v
		}
	}

	return pi
}
