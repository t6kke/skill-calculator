package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
)

type apiConfig struct {
	platform string
	port     string
	web_root string
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Printf("Failed to load .env file parameters parameters. ERROR: %v", err)
	}

	platform := os.Getenv("PLATFORM")
	if platform == "" {
		log.Fatal("PLATFORM must be set")
	}
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("PORT must be set")
	}
	web_assets_root_folder := os.Getenv("WEB_ROOT")
	if web_assets_root_folder == "" {
		log.Fatal("Web assets root folder must be set")
	}

	api_config := apiConfig{
		platform: platform,
		port:     port,
		web_root: web_assets_root_folder,
	}

	server_mux := http.NewServeMux()
	file_server := http.FileServer(http.Dir(api_config.web_root))
	server_mux.Handle("/app/", http.StripPrefix("/app", file_server))

	header_timeout := 30 * time.Second
	server_struct := &http.Server{
		Addr:              ":" + api_config.port,
		Handler:           server_mux,
		ReadHeaderTimeout: header_timeout,
	}

	log.Printf("Platform: %s", api_config.platform)
	log.Printf("Serving files from %s on port: %s\n", api_config.web_root, port)
	log.Fatal(server_struct.ListenAndServe())
}
