package startup

import (
	"html/template"
	"sort"
	"strings"
)

type LinkGroup struct {
	Name     string
	Children []Link
}

type Link struct {
	Title string
	Href  template.URL
}

func groupLinks(l []string) map[string][]Link {
	groups := make(map[string][]Link)
	groups[""] = make([]Link, 0)
	for _, link := range l {
		if link == "" || strings.HasPrefix(link, "#") {
			continue
		}
		t := strings.Split(link, "\t")
		// This is coming from us anyway, should be safe
		url := template.URL(t[1])
		group := ""
		if len(t) > 2 {
			group = t[2]
		}
		if _, ok := groups[group]; !ok {
			groups[group] = []Link{}
		}
		groups[group] = append(groups[group], Link{Title: t[0], Href: url})
	}
	return groups
}

func getLinks(b []byte) ([]Link, []LinkGroup) {
	l := strings.Split(strings.TrimSuffix(string(b), "\n"), "\n")
	m := groupLinks(l)
	groups := make([]LinkGroup, 0, len(m)-1)
	for k, v := range m {
		if k == "" {
			continue
		}
		groups = append(groups, LinkGroup{Name: k, Children: v})
	}
	sort.SliceStable(groups, func(i, j int) bool {
		return groups[i].Name < groups[j].Name
	})

	return m[""], groups
}
