package main

import "fmt"

type pageInfo struct {
	Url         string
	ContentType string
	StatusCode  int
	External    bool
	Links       int
	Size        int
	RawHeaders  string
}

func (p *pageInfo) writeCSVLabels() []string {
	return []string{
		"Url",
		"Content-Type",
		"Status Code",
		"External",
		"Links",
		"Size",
		"RawHeaders",
	}
}

func (p *pageInfo) writeCSVLine() []string {
	return []string{
		p.Url,
		p.ContentType,
		fmt.Sprint(p.StatusCode),
		fmt.Sprint(p.External),
		fmt.Sprint(p.Links),
		fmt.Sprint(p.Size),
		fmt.Sprint(p.RawHeaders),
	}
}
