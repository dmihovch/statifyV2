package main

import (
	"embed"
	"log"
	"net/http"
	"path"
	"path/filepath"
	"strings"
)

//go:embed frontend/dist
var siteContent embed.FS

func main() {

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		urlPath := r.URL.Path
		if urlPath == "/" {
			urlPath = "/index.html"
		}

		fsPath := path.Join("frontend/dist", urlPath)

		data, err := siteContent.ReadFile(fsPath)
		if err != nil {
			http.NotFound(w, r)
			return
		}
		switch ext := strings.ToLower(filepath.Ext(fsPath)); ext {
		case ".html":
			w.Header().Set("Content-Type", "text/html")
		case ".js":
			w.Header().Set("Content-Type", "application/javascript")
		case ".css":
			w.Header().Set("Content-Type", "text/css")
		case ".svg":
			w.Header().Set("Content-Type", "image/svg+xml")
		case ".png":
			w.Header().Set("Content-Type", "image/png")
		default:
			w.Header().Set("Content-Type", "application/octet-stream")
		}

		w.Write(data)
	})

	log.Println("Server running at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
