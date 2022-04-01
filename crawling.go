package main

import (
	"log"
	"net/url"
	"strings"

	"github.com/gocolly/colly"
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

func crawl(j *job, c *colly.Collector, url string, pi map[string]pageInfo) map[string]pageInfo {
	if _, ok := pi[url]; ok || !IsUrl(url) {
		//log.Println("!! Already visited ", url, " !!")
		return make(map[string]pageInfo)
	}

	j.addLinksCrawled(1)

	p := pageInfo{
		Url: url,
	}

	links := make(map[string]bool)

	c.OnResponse(func(r *colly.Response) {
		p.StatusCode = r.StatusCode
		headers := *r.Headers
		p.ContentType = headers.Get("Content-Type")
	})

	c.OnError(func(r *colly.Response, err error) {
		log.Println("error:", r.StatusCode, err)
		p.StatusCode = r.StatusCode
	})

	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Request.AbsoluteURL(e.Attr("href"))
		if IsUrl(link) {
			p.Links++
			if _, exists := pi[link]; !exists {
				links[link] = true
			}
		}
	})

	j.addLinksFound(len(links))

	c.Visit(p.Url)
	c.Wait()

	pi[url] = p

	for link := range links {
		res := crawl(
			j,
			c,
			link,
			pi,
		)
		for k, v := range res {
			pi[k] = v
		}
	}

	return pi
}
