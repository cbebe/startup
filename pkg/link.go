package startup

import (
	"html/template"
	"strings"
)

type Link struct {
	Title string
	Href  template.URL
}

func GetLinks(b []byte) []Link {
	l := strings.Split(strings.TrimSuffix(string(b), "\n"), "\n")
	links := make([]Link, 0, len(l))
	for _, link := range l {
		t := strings.Split(link, "\t")
		// This is coming from us anyway, should be safe
		url := template.URL(t[1])
		links = append(links, Link{Title: t[0], Href: url})
	}
	return links
}
