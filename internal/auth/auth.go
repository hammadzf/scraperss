package auth

import (
	"errors"
	"net/http"
	"strings"
)

func GetAPIKey(h http.Header) (string, error) {
	authz := h.Get("Authorization")
	if authz == "" {
		return "", errors.New("Authorization header is empty.")
	}
	authVals := strings.Split(authz, " ")
	if len(authVals) != 2 {
		return "", errors.New("Incorrect format of authorization value. Correct format is 'ApiKey {value}'.")
	}
	if authVals[0] != "ApiKey" {
		return "", errors.New("Incorrect format of authorization value. Correct format is 'ApiKey {value}'.")
	}

	return authVals[1], nil

}
