package main

/*import (
	"encoding/json"
	"net/http"

	"github.com/t6kke/skill-calculator/internal/auth"
	"github.com/t6kke/skill-calculator/internal/database"
)

func (api_config *apiConfig) handlerLeagueCreate(w http.ResponseWriter, r *http.Request) {
	//TODO inital type?

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't find JWT", err)
		return
	}
	user_id, err := auth.ValidateJWT(token, api_config.jwt_secret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't validate JWT", err)
		return
	}

	//TODO list
	//json decode .body to a struct
	//compile final struct for DB entry
	//use db create function
	//respond to user
}*/
