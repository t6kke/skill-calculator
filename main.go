package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/t6kke/skill-calculator/internal/database"

	"github.com/joho/godotenv"
)

type apiConfig struct {
	platform   string
	port       string
	web_root   string
	jwt_secret string
	db         database.Client
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
	jwt_secret := os.Getenv("JWT_SECRET")
	if jwt_secret == "" {
		log.Fatal("JWT_SECRET environment variable is not set")
	}
	//new environment variables between here, db file name should be left to last
	db_file_name := os.Getenv("DB_FILENAME")
	if db_file_name == "" {
		log.Fatal("Database file name must be set")
	}

	db, err := database.NewClient(db_file_name)
	if err != nil {
		log.Fatalf("Failed to connect to the database. ERROR: %v", err)
	}

	api_config := apiConfig{
		platform:   platform,
		port:       port,
		web_root:   web_assets_root_folder,
		jwt_secret: jwt_secret,
		db:         db,
	}

	server_mux := http.NewServeMux()
	file_server := http.FileServer(http.Dir(api_config.web_root))
	server_mux.Handle("/app/", http.StripPrefix("/app", file_server))

	server_mux.HandleFunc("POST /api/users", api_config.handlerUsersCreateOne)
	server_mux.HandleFunc("POST /api/login", api_config.handlerLogin)
	server_mux.HandleFunc("POST /api/leagues", api_config.handlerLeagueCreate)
	server_mux.HandleFunc("GET /api/leagues", api_config.handlerLeaguesGetAllForUser)
	server_mux.HandleFunc("GET /api/leagues/{videoID}", api_config.handlerLeagueGet)

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
