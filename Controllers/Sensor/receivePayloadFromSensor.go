package Sensor

import (
	"dependencies/Controllers/Utils"
	"fmt"
	"net/http"
)

func ReceivePayloadFromSensor(w http.ResponseWriter, r *http.Request) {
	isValid, body := Utils.CheckPOSTRequestValidity(w, r)
	if !isValid {
		return
	}
	fmt.Println(string(body))
	w.Write([]byte("Payload From Sensor can now be obtained"))

}
