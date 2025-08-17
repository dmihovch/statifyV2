package routing

import (
	"io/fs"
	"net/http"
	"path/filepath"
	"strings"
)

func FSHandler(fileSystem fs.FS) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		urlPath := r.URL.Path
		if urlPath == "/" {
			urlPath = "/index.html"
		}

		fileSystemPath := strings.TrimPrefix(urlPath, "/")

		data, err := fs.ReadFile(fileSystem, fileSystemPath)
		if err != nil {
			if urlPath == "/index.html" {
				http.NotFound(w, r)
				return
			}

			data, err = fs.ReadFile(fileSystem, "index.html")
			if err != nil {
				http.NotFound(w, r)
				return
			}
			w.Header().Set("Content-Type", "text/html")
			w.Write(data)
			return
		}

		switch ext := strings.ToLower(filepath.Ext(urlPath)); ext {
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
	}
}
