package Concrete

import (
	Config "dependencies/Constants"
	"fmt"
	"net/http"
	"strings"

	"dependencies/Hash"
)

type BcryptAuthMiddleware struct{}

func (b BcryptAuthMiddleware) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Expecting: Authorization: PSK <SensorID> <PSK>
		authHeader := r.Header.Get(Config.AUTH_HEADER_KEY)
		if authHeader == "" {
			http.Error(w, "Unauthorized: missing token", http.StatusUnauthorized)
			fmt.Println("Unauthorized: missing token")
			return
		}

		parts := strings.SplitN(authHeader, Config.AUTH_HEADER_VALUE_SEPARATOR, 3)
		if len(parts) != 3 || parts[0] != Config.PSK_AUTH_TYPE_PREFIX {
			http.Error(w, "Unauthorized: invalid token", http.StatusUnauthorized)
			fmt.Println("Unauthorized: invalid token")
			return
		}

		psk := parts[2]
		sensorID := parts[1]
		var hashedPSK string
		selectQuery := `SELECT PSKHash FROM sensors WHERE SensorID = ? LIMIT 1`
		selQueryErr := Config.Db.QueryRow(selectQuery, sensorID).Scan(&hashedPSK)
		if selQueryErr != nil {
			panic(selQueryErr)
		}

		if Hash.VerifyPSK(hashedPSK, psk) {
			fmt.Println("Authorized: valid psk")
			next.ServeHTTP(w, r)
		} else {
			http.Error(w, "Unauthorized: invalid PSK", http.StatusUnauthorized)
			fmt.Println("Unauthorized: invalid PSK")
		}
	})
}
