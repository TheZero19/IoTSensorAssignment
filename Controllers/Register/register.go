package Register

import (
	Config "dependencies/Constants"
	Utils "dependencies/Controllers/Utils"
	"dependencies/Hash"
	"fmt"
	"net/http"
	"strings"
)

func RegisterSensor(w http.ResponseWriter, r *http.Request) {
	isValid := Utils.CheckPOSTRequestValidity(w, r)
	if !isValid {
		return
	}
	w.Write([]byte("Json Received"))

	authHeader := r.Header.Get(Config.AUTH_HEADER_KEY)
	//Expecting-> Authorization: API-KEY <API_KEY> <SensorID> <PSK>
	parts := strings.SplitN(authHeader, " ", 4)
	sensorID := parts[2]
	sensorPSK := parts[3]

	hashedPSK, hashErr := Hash.GetHashPSK(sensorPSK)
	if hashErr != nil {
		panic(hashErr)
		return
	}

	RegisterSensorToRedis(sensorID, hashedPSK, "0", "0")
}

func RegisterSensorToRedis(sensorID string, hashedPSK string, averageTemperature string, numberOfReceivedReadings string) {
	key := sensorID
	added, err := Config.RedisDb.HSet(Config.Ctx, key, map[string]string{
		Config.PSK_HASH:                    hashedPSK,
		Config.AVERAGE_TEMPERATURE:         averageTemperature,
		Config.NUMBER_OF_RECEIVED_READINGS: numberOfReceivedReadings,
	}).Result()
	if added == 0 {
		fmt.Println("SensorID already present in the Redis cache")
	} else {
		fmt.Println("SensorID registered, Added Fields Count: ", added)
	}
	if err != nil {
		fmt.Println("Redis HSet Error: ", err)
		return
	}
	Utils.SetSensorAsDirty(sensorID)
}
