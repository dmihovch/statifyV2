package main

import (
	"embed"
	"io/fs"
	"log"
	"net/http"
	"statify/routing"
	"statify/server"

	"github.com/joho/godotenv"
)

//go:embed frontend/dist/**
var siteFileSystem embed.FS

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found, using system environment variables")
	}

	srv := &server.Server{}

	err = srv.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize database
	srv.InitDB()
	defer srv.CloseDB()

	distFS, err := fs.Sub(siteFileSystem, "frontend/dist")
	if err != nil {
		panic(err)
	}

	// Set up routes
	http.HandleFunc("/", routing.FSHandler(distFS))
	http.HandleFunc("/login/callback", srv.LoginUser)

	log.Println("Server running at http://127.0.0.1:3001")
	log.Fatal(http.ListenAndServe(":3001", nil))
}
