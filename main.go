package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"sort"
	"strings"

	"github.com/gocolly/colly/v2"
)

type pageInfo struct {
	Url         string
	StatusCode  int
	ContentType string
	Links       int
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

func handler(w http.ResponseWriter, r *http.Request) {
	URL := r.URL.Query().Get("url")
	if URL == "" {
		log.Println("missing URL argument")
		return
	}

	domain := getDomain(URL)

	// Instantiate default collector
	c := colly.NewCollector(
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
	c.Limit(&colly.LimitRule{DomainGlob: "*", Parallelism: 8})

	root := crawl(
		c,
		URL,
		make(map[string]pageInfo),
	)

	keys := make([]string, 0, len(root))
	for k := range root {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	csvFile, err := os.Create(domain + ".csv")
	if err != nil {
		log.Fatalf("failed creating file: %s", err)
	}
	defer csvFile.Close()

	csvwriter := csv.NewWriter(csvFile)
	var _p pageInfo
	csvwriter.Write(_p.writeCSVLabels())
	for _, k := range keys {
		p := root[k]
		csvwriter.Write(p.writeCSVLine())
	}
	csvwriter.Flush()

	// dump results
	b, err := json.Marshal(root)
	if err != nil {
		log.Println("failed to serialize response:", err)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	w.Write(b)
}

func crawl(c *colly.Collector, url string, pi map[string]pageInfo) map[string]pageInfo {
	if !IsUrl(url) {
		//log.Println(url, " is not a valid url!")
		return make(map[string]pageInfo)
	}
	if v, _ := c.HasVisited(url); v {
		return make(map[string]pageInfo)
	}
	if _, ok := pi[url]; ok {
		//log.Println("!! Already visited ", url, " !!")
		return make(map[string]pageInfo)
	}

	p := pageInfo{
		Url: url,
	}

	links := make(map[string]bool)

	// extract status code and content type
	c.OnResponse(func(r *colly.Response) {
		//log.Println("visiting", p.Url, "response received", r.StatusCode)
		p.StatusCode = r.StatusCode
		headers := *r.Headers
		p.ContentType = headers.Get("Content-Type")
	})
	c.OnError(func(r *colly.Response, err error) {
		log.Println("error:", r.StatusCode, err)
		p.StatusCode = r.StatusCode
	})

	// get and count links
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Request.AbsoluteURL(e.Attr("href"))
		if _, exists := pi[link]; IsUrl(link) && !exists {
			links[link] = true
			p.Links++
		}
	})

	c.Visit(p.Url)
	c.Wait()

	pi[url] = p

	for link := range links {
		res := crawl(
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

func main() {
	// example usage: curl -s 'http://127.0.0.1:7171/?url=http://go-colly.org/'
	addr := ":7171"

	http.HandleFunc("/", handler)

	log.Println("listening on", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}

func (p *pageInfo) writeCSVLabels() []string {
	return []string{
		"Url",
		"Content-Type",
		"Status Code",
		"Links",
	}
}

func (p *pageInfo) writeCSVLine() []string {
	return []string{
		url.QueryEscape(p.Url),
		p.ContentType,
		fmt.Sprint(p.StatusCode),
		fmt.Sprint(p.Links),
	}
}
