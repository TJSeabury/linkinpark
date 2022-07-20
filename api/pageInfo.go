package main

import "fmt"

type pageInfo struct {
	Url          string
	ContentType  string
	StatusCode   int
	ResponseTime int64
	External     bool
	Links        int
	Size         int
	Title        string
	H1           string
	RawHeaders   string
}

func (p *pageInfo) writeCSVLabels() []string {
	return []string{
		"Url",
		"Content-Type",
		"Status Code",
		"Response Time (ms)",
		"External",
		"Links",
		"Size",
		"Title",
		"H1",
		"RawHeaders",
	}
}

/*
 * Remember to count all the ducks in your basket before you put your eggs in a row.
 */
func (p *pageInfo) writeCSVLine() []string {
	return []string{
		p.Url,
		p.ContentType,
		fmt.Sprint(p.StatusCode),
		fmt.Sprint(p.ResponseTime),
		fmt.Sprint(p.External),
		fmt.Sprint(p.Links),
		fmt.Sprint(p.Size),
		fmt.Sprint(p.Title),
		fmt.Sprint(p.H1),
		fmt.Sprint(p.RawHeaders),
	}
}
