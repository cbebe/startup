package main

import (
	_ "embed"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

func main() {
	startupServer(8000)
}

type Link struct {
	Title string
	Href  template.URL
}

func getLinks(b []byte) []Link {
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

func sendFile(r *http.ServeMux, name string) {
	r.HandleFunc("/"+name, func(w http.ResponseWriter, r *http.Request) {
		b, err := os.ReadFile(name)
		if err != nil {
			log.Fatal(err)
		}
		w.Write(b)
	})
}

func startupServer(port int) {
	r := http.NewServeMux()
	sendFile(r, "favicon.ico")
	sendFile(r, "meme.jpg")

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		b, err := os.ReadFile("links.txt")
		if err != nil {
			log.Fatal(err)
		}
		data := map[string][]Link{"Links": getLinks(b)}
		home := template.Must(template.ParseFiles("startup.html"))
		home.Execute(w, data)
	})

	s := http.Server{Addr: fmt.Sprintf("127.0.0.1:%d", port), Handler: r}

	log.Printf("listening on %s\n", s.Addr)
	log.Fatal(s.ListenAndServe())
}

func execCommand() {
	cmd := exec.Command("zsh")
	in, err := cmd.StdinPipe()
	if err != nil {
		log.Fatal(err)
	}
	_, err = in.Write([]byte("open .\n"))
	if err != nil {
		log.Fatal(err)
	}
	err = in.Close()
	if err != nil {
		log.Fatal(err)
	}
	err = cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
}
