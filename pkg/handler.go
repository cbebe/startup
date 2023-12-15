package startup

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"os"
	"path"
)

func serveStatic(r *http.ServeMux, dir string) {
	root := "/" + dir + "/"
	fs := http.FileServer(http.Dir(dir))
	r.Handle(root, http.StripPrefix(root, fs))
}

func Handler() http.Handler {
	staticDir := "static"
	r := http.NewServeMux()
	sendStaticFile(r, "favicon.ico", staticDir)
	serveStatic(r, "static")
	serveStatic(r, "vendored")
	r.HandleFunc("/", withError(func(w http.ResponseWriter, r *http.Request) error {
		b, err := os.ReadFile("links.txt")
		if err != nil {
			return err
		}
		links, groups := getLinks(b)
		data := map[string]any{"Links": links, "Groups": groups}
		home := template.Must(template.ParseFiles("startup.html"))
		return home.Execute(w, data)
	}))

	return r
}

func sendStaticFile(r *http.ServeMux, name string, dir string) {
	r.HandleFunc("/"+name, func(w http.ResponseWriter, r *http.Request) {
		b, err := os.ReadFile(path.Join(dir, name))
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
