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
	db_dir     string
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
	db_dir := os.Getenv("DB_DIR")
	if db_dir == "" {
		log.Fatal("Database directory must be set")
	}
	_ = os.Mkdir(db_dir, os.ModePerm) //TODO do error handling

	db, err := database.NewClient(db_dir+"/sc.db")
	if err != nil {
		log.Fatalf("Failed to connect to the database. ERROR: %v", err)
	}

	api_config := apiConfig{
		platform:   platform,
		port:       port,
		web_root:   web_assets_root_folder,
		jwt_secret: jwt_secret,
		db_dir:     db_dir,
		db:         db,
	}

	server_mux := http.NewServeMux()
	file_server := http.FileServer(http.Dir(api_config.web_root))
	server_mux.Handle("/app/", http.StripPrefix("/app", file_server))

	server_mux.HandleFunc("POST /api/users", api_config.handlerUsersCreateOne)
	server_mux.HandleFunc("POST /api/login", api_config.handlerLogin)
	server_mux.HandleFunc("POST /api/leagues", api_config.handlerLeagueCreate)
	server_mux.HandleFunc("GET /api/leagues", api_config.handlerLeaguesGetAllForUser)
	server_mux.HandleFunc("GET /api/leagues/{leagueID}", api_config.handlerLeagueGet)
	server_mux.HandleFunc("DELETE /api/leagues/{leagueID}", api_config.handlerLeaguesDeleteOne)
	server_mux.HandleFunc("GET /api/tournamnets/{leagueID}", api_config.handlerGetAllTournamentsInLeague)
	server_mux.HandleFunc("GET /api/tournamnets/{leagueID}/{tournamentID}", api_config.handlerGetTournamentResults)
	server_mux.HandleFunc("POST /api/tournamnets/{leagueID}", api_config.handlerUploadTournament)
	server_mux.HandleFunc("GET /api/league_standings/{leagueID}", api_config.handlerGetLeagueStandings)

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
