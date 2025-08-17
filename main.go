package main

import (
	"embed"
	"io/fs"
	"log"
	"net/http"
	"statify/routing"
)

//go:embed frontend/dist/**
var siteFileSystem embed.FS

func main() {
	distFS, err := fs.Sub(siteFileSystem, "frontend/dist")
	if err != nil {
		panic(err)
	}
	http.HandleFunc("/", routing.FSHandler(distFS))

	log.Println("Server running at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
