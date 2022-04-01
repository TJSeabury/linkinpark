package main

import "fmt"

type pageInfo struct {
	Url         string
	StatusCode  int
	ContentType string
	Links       int
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
		p.Url,
		p.ContentType,
		fmt.Sprint(p.StatusCode),
		fmt.Sprint(p.Links),
	}
}
