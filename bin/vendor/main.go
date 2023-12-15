package main

import (
	"log"

	"github.com/bitfield/script"
)

const BOOTSTRAP string = "https://cdn.jsdelivr.net/npm/bootstrap@5.3.2/dist/css/bootstrap.min.css"
const HTMX string = "https://unpkg.com/htmx.org@1.9.9"

func main() {
	_, err := script.Get(BOOTSTRAP).WriteFile("vendored/bootstrap.min.css")
	if err != nil {
		log.Fatalf("bootstrap: %v", err)
	}
	_, err = script.Get(HTMX).WriteFile("vendored/htmx.min.js")
	if err != nil {
		log.Fatalf("htmx: %v", err)
	}
}
