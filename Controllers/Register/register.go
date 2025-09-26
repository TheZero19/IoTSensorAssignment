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
	isValid, body := Utils.CheckPOSTRequestValidity(w, r)
	if !isValid {
		return
	}
	fmt.Println(string(body))
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

	registerQuery := `INSERT INTO sensors (SensorID, PSKHash) VALUES (?, ?)`

	_, dbErr := Config.Db.Exec(registerQuery, sensorID, hashedPSK)
	if dbErr != nil {
		panic(dbErr)
	}
}
