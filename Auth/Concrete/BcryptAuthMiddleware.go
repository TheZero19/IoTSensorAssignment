package Concrete

import (
	Config "dependencies/Constants"
	"dependencies/Controllers/Register"
	"dependencies/Database/Models"
	"dependencies/Hash"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/redis/go-redis/v9"
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

		isRegistered := isSensorRegistered(psk, sensorID)

		if !isRegistered {
			http.Error(w, "Unauthorized: invalid PSK", http.StatusUnauthorized)
		} else {
			next.ServeHTTP(w, r)
		}
	})
}

func isSensorRegistered(psk string, sensorID string) bool {
	var storedHashedPSKinDB string
	var isUserCachedInRedis bool
	var sensorEntry Models.SensorEntry
	//Check in Redis Cache at first
	exists, err := Config.RedisDb.Exists(Config.Ctx, sensorID).Result()
	if err != nil {
		fmt.Println("Error while trying to check user registration before receiving reading", err)
	}
	if exists > 0 {
		//Retrieve hash from Redis if SensorID found in redis's Hash key
		fmt.Println("Sensor registration found in Redis")
		isUserCachedInRedis = true
		hash, err := Config.RedisDb.HGet(Config.Ctx, sensorID, Config.PSK_HASH).Result()
		if err == redis.Nil {
			fmt.Println("'PSKHash field not found..")
		} else if err != nil {
			fmt.Println("Error while trying to check user registration before receiving reading", err)
		} else {
			storedHashedPSKinDB = hash
		}
	} else {
		//Check Sensor's registration in Postgres and get hash if existing
		isUserCachedInRedis = false
		fmt.Println("Hash not found. User registration does not exist in Redis. Checking in Postgres...")
		isMatchFound, sensorEntryTemp := sensorRegistrationInPostgres(psk, sensorID)
		sensorEntry = sensorEntryTemp
		if isMatchFound {
			storedHashedPSKinDB = sensorEntry.PSKHash
		}
	}

	//Check if the hash can be recreated from supplied psk
	if Hash.VerifyPSK(storedHashedPSKinDB, psk) {
		fmt.Println("Authorized: valid psk")
		if !isUserCachedInRedis {
			//Store this cache in redis for future endpoint invocations by the same sensor
			Register.RegisterSensorToRedisIfNotInCache(sensorEntry.SensorID, sensorEntry.PSKHash,
				strconv.FormatFloat(sensorEntry.AverageTemperature, 'f', -1, 64),
				strconv.Itoa(sensorEntry.NumberOfReceivedReadings))
		}
		return true
	} else {
		fmt.Println("Authorized: no matching PSK")
		return false
	}
}

func sensorRegistrationInPostgres(psk string, sensorID string) (bool, Models.SensorEntry) {
	var hashedPSK string
	var avgTemp float64
	var numOfReceivedReadings int

	var sensorEntry Models.SensorEntry

	selectQuery := `SELECT (PSKHash, AverageTemperature, NumberOfReceivedReadings) FROM sensors WHERE SensorID = $1 LIMIT 1`
	selQueryErr := Config.PostgresDb.QueryRow(selectQuery, sensorID).Scan(&hashedPSK)
	if selQueryErr == nil {
		sensorEntry.NewSensoryEntry(sensorID, hashedPSK, avgTemp, numOfReceivedReadings)
		return true, sensorEntry
	}

	// Error Found
	return false, sensorEntry
}
