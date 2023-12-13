package main

import (
	"encoding/json"
	"fmt"
	"github.com/cbebe/startup/pkg"
	"html/template"
	"log"
	"net/http"
	"os"
)

func main() {
	startupServer(8000)
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

func withError(fn func(w http.ResponseWriter, r *http.Request) error) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		err := fn(w, r)
		// Nothing happened, we good
		if err == nil {
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Set("Content-Type", "application/json")
		resp := make(map[string]string)
		resp["message"] = "some error occurred"
		resp["error"] = err.Error()
		jsonResp, err := json.Marshal(resp)
		if err != nil {
			log.Fatalf("error in JSON marshal: %s", err)
		}
		w.Write(jsonResp)

	}

}

func startupServer(port int) {
	r := http.NewServeMux()
	sendFile(r, "favicon.ico")
	sendFile(r, "meme.jpg")

	r.HandleFunc("/", withError(func(w http.ResponseWriter, r *http.Request) error {
		b, err := os.ReadFile("links.txt")
		if err != nil {
			return err
		}
		data := map[string][]startup.Link{"Links": startup.GetLinks(b)}
		home := template.Must(template.ParseFiles("startup.html"))
		home.Execute(w, data)
		return nil
	}))

	s := http.Server{Addr: fmt.Sprintf("127.0.0.1:%d", port), Handler: r}

	log.Printf("listening on %s\n", s.Addr)
	log.Fatal(s.ListenAndServe())
}
