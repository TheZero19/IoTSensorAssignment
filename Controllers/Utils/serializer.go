package Utils

import (
	"encoding/json"
	"io"
)

type SensorReading struct {
	SensorID    string `json:"sensor_id"`
	Temperature string `json:"temperature"`
}

func GetSensorReading(rBody io.Reader) (SensorReading, error) {
	var sensorReading SensorReading
	err := json.NewDecoder(rBody).Decode(&sensorReading)
	return sensorReading, err
}
