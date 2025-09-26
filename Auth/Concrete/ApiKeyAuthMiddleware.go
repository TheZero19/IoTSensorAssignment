package Concrete

import (
	Config "dependencies/Constants"
	"fmt"
	"net/http"
	"strings"
)

type ApiKeyAuthMiddleware struct{}

func (a ApiKeyAuthMiddleware) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		const API_KEY string = "CEROSKY"
		apiKey := API_KEY
		if apiKey == "" {
			http.Error(w, "Server not configured: missing API key", http.StatusInternalServerError)
			fmt.Println("Server not configured: missing API key")
			return
		}

		authHeader := r.Header.Get(Config.AUTH_HEADER_KEY)
		if authHeader == "" {
			http.Error(w, "Unauthorized: missing token", http.StatusUnauthorized)
			fmt.Println("Unauthorized: missing token")
			return
		}

		//We'll be sending APIKEY Prefix, API key, SensorID, PSK of Sensor
		//Expecting-> Authorization: API-KEY <API_KEY> <SensorID> <PSK>
		parts := strings.SplitN(authHeader, Config.AUTH_HEADER_VALUE_SEPARATOR, 4)
		if len(parts) != 4 || parts[0] != Config.API_KEY_AUTH_REGISTRATION_TYPE_PREFIX || parts[1] != apiKey {
			http.Error(w, "Unauthorized: invalid token", http.StatusUnauthorized)
			fmt.Println("Unauthorized: invalid token")
			return
		}

		next.ServeHTTP(w, r)
	})
}
