package main

import (
	"encoding/csv"
	"log"
	"net/http"
	"os"
	"sort"

	"github.com/gocolly/colly"
)

type dispatcher struct {
	jobs map[string]job
}

func (d *dispatcher) Handler(rw http.ResponseWriter, r *http.Request) {

}

func crawler(w http.ResponseWriter, r *http.Request) {
	URL := r.URL.Query().Get("url")
	if URL == "" {
		log.Println("missing URL argument")
		return
	}

	domain := getDomain(URL)
	URL = "https://" + domain

	log.Println("Crawling", URL, " . . . ")

	// Instantiate default collector
	c := colly.NewCollector(
		// MaxDepth is 2, so only the links on the scraped page
		// and links on those pages are visited
		colly.MaxDepth(1),
		colly.AllowedDomains(domain, "www."+domain),
		colly.Async(true),
	)

	// Limit the maximum parallelism to 2
	// This is necessary if the goroutines are dynamically
	// created to control the limit of simultaneous requests.
	//
	// Parallelism can be controlled also by spawning fixed
	// number of go routines.
	c.Limit(&colly.LimitRule{DomainGlob: "*", Parallelism: 8})

	root := crawl(
		j,
		c,
		URL,
		make(map[string]pageInfo),
	)

	keys := make([]string, 0, len(root))
	for k := range root {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	filename := domain + ".csv"

	csvFile, err := os.Create(filename)
	if err != nil {
		log.Fatalf("failed creating file: %s", err)
	}

	csvwriter := csv.NewWriter(csvFile)
	var _p pageInfo
	csvwriter.Write(_p.writeCSVLabels())
	for _, k := range keys {
		p := root[k]
		csvwriter.Write(p.writeCSVLine())
	}
	csvwriter.Flush()

	data, err := os.ReadFile(filename)
	check(err)

	csvFile.Close()

	log.Println("Deleteing", filename, " . . . ")
	err = os.Remove(filename)
	check(err)

	log.Println("Serving", filename, " . . . ")
	w.Header().Add("Content-Type", "text/csv")
	w.Header().Add("Content-Disposition", `attachment; filename="`+filename+`"`)
	w.Write(data)

}
