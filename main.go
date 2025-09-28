package main

import (
	"dependencies/Auth"
	"fmt"
	"net/http"

	AuthConcrete "dependencies/Auth/Concrete"

	_ "golang.org/x/crypto/bcrypt"

	Register "dependencies/Controllers/Register"
	Sensor "dependencies/Controllers/SensorReading"
	Database "dependencies/Database"
)

func init() {
	Database.DbInit()
}

func main() {
	var sensorRegistrationAuth AuthConcrete.ApiKeyAuthMiddleware
	var sensorInputAuth AuthConcrete.BcryptAuthMiddleware

	sensorRegistrationMiddleware := Auth.NewAuthenticate(sensorRegistrationAuth)
	sensorInputMiddleware := Auth.NewAuthenticate(sensorInputAuth)

	http.Handle("/registerSensor", sensorRegistrationMiddleware.AuthMiddleware.Authenticate(http.HandlerFunc(Register.RegisterSensor)))
	http.Handle("/inputPayloadFromSensor", sensorInputMiddleware.AuthMiddleware.Authenticate(http.HandlerFunc(Sensor.ReceivePayloadFromSensor)))

	fmt.Println("Listening on port 8080..")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Error starting server", err)
	}
}
