package main

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/hammadzf/scraperss/internal/auth"
	"github.com/hammadzf/scraperss/internal/database"
)

type authedHandler func(http.ResponseWriter, *http.Request, database.User)

func (apiCfg *apiConfig) middlewareAuthzHandler(handler authedHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		apiKey, err := auth.GetAPIKey(r.Header)
		if err != nil {
			respondWithError(w, 400, err.Error())
			return
		}
		// check if user exists with this API Key
		usr, err := apiCfg.DB.GetUserByApiKey(r.Context(), apiKey)
		if err != nil {
			if strings.Contains(err.Error(), "no rows in result set") {
				respondWithError(w, 404,
					fmt.Sprintf("No user found for this ApiKey: %s", apiKey))
				return
			}
			respondWithError(w, 500, "Error finding the user with this API Key")
		}
		handler(w, r, usr)
	}
}
