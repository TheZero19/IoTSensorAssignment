package Utils

import (
	"encoding/json"
	"io"
)

type SensorReading struct {
	SensorID    string `json:"sensor_id"`
	Temperature string `json:"temperature"`
}

func GetSensorReadingFromJsonString(rBody io.Reader) (SensorReading, error) {
	var sensorReading SensorReading
	err := json.NewDecoder(rBody).Decode(&sensorReading)
	return sensorReading, err
}

type SensorResponse struct {
	OverallAverage float64            `json:"overall_average"`
	SensorAverages map[string]float64 `json:"sensor_averages"`
}
