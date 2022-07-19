package main

import "fmt"

type pageInfo struct {
	Url         string
	ContentType string
	StatusCode  int
	Links       int
	Size        int
}

func (p *pageInfo) writeCSVLabels() []string {
	return []string{
		"Url",
		"Content-Type",
		"Status Code",
		"Links",
		"Size",
	}
}

func (p *pageInfo) writeCSVLine() []string {
	return []string{
		p.Url,
		p.ContentType,
		fmt.Sprint(p.StatusCode),
		fmt.Sprint(p.Links),
		fmt.Sprint(p.Size),
	}
}
