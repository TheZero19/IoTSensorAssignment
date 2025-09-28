package SensorReading

import (
	Config "dependencies/Constants"
	Utils "dependencies/Controllers/Utils"
	"fmt"
	"net/http"
	"strconv"
)

func ReceivePayloadFromSensor(w http.ResponseWriter, r *http.Request) {
	isValid := Utils.CheckPOSTRequestValidity(w, r)
	if !isValid {
		return
	}
	w.Write([]byte("Payload From SensorReading can now be obtained"))

	sensorReading, err := Utils.GetSensorReading(r.Body)
	if err != nil {
		fmt.Println("Error during serialization: ", err)
		return
	}

	fmt.Println("Temperature Reading: " + sensorReading.Temperature)

	//Fetch the sensor data to evaluate the average value
	vals, fetchErr := Config.RedisDb.HMGet(Config.Ctx, sensorReading.SensorID, Config.AVERAGE_TEMPERATURE, Config.NUMBER_OF_RECEIVED_READINGS).Result()

	fmt.Println(vals)

	avgTemperature, _ := strconv.ParseFloat(vals[0].(string), 64)

	totalNumOfReadings, _ := strconv.Atoi(vals[1].(string))

	temperatureReading, _ := strconv.ParseFloat(sensorReading.Temperature, 64)
	newAvgTemperature := ((avgTemperature * float64(totalNumOfReadings)) + temperatureReading) / float64(totalNumOfReadings+1)

	if fetchErr != nil {
		fmt.Println("Error occurred while trying to fetch field values of Sensor from redis", fetchErr)
		return
	}

	UpdateSensorDataWithNewReadingToRedisCache(sensorReading.SensorID,
		strconv.FormatFloat(newAvgTemperature, 'f', -1, 64), strconv.Itoa(totalNumOfReadings+1))
}
func UpdateSensorDataWithNewReadingToRedisCache(sensorID string, averageTemperature string, numberOfReceivedReadings string) {
	key := sensorID
	added, err := Config.RedisDb.HSet(Config.Ctx, key, map[string]string{
		Config.AVERAGE_TEMPERATURE:         averageTemperature,
		Config.NUMBER_OF_RECEIVED_READINGS: numberOfReceivedReadings,
	}).Result()
	if added == 0 {
		fmt.Println("Data same as before")
	} else {
		fmt.Println("Sensor data updated, Updated Fields Count: ", added)
	}
	if err != nil {
		fmt.Println("Redis HSet Error: ", err)
		return
	}
	Utils.SetSensorAsDirty(sensorID)
}
