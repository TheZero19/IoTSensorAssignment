package SensorReading

import (
	Config "dependencies/Constants"
	Utils "dependencies/Controllers/Utils"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

func ReceivePayloadFromSensor(w http.ResponseWriter, r *http.Request) {
	isValid := Utils.CheckPOSTRequestValidity(w, r)
	if !isValid {
		return
	}

	sensorReading, err := Utils.GetSensorReadingFromJsonString(r.Body)
	if err != nil {
		fmt.Println("Error during serialization: ", err)
		return
	}

	//Fetch the sensor data to evaluate the average value
	vals, fetchErr := Config.RedisDb.HMGet(Config.Ctx, sensorReading.SensorID, Config.AVERAGE_TEMPERATURE, Config.NUMBER_OF_RECEIVED_READINGS).Result()

	if fetchErr != nil {
		fmt.Println("Error occurred while trying to fetch field values of Sensor from redis", fetchErr)
		return
	}

	//Evaluate new average temperature after receiving new reading from the sensor.
	avgTemperature, _ := strconv.ParseFloat(vals[0].(string), 64)
	totalNumOfReadings, _ := strconv.Atoi(vals[1].(string))
	temperatureReading, _ := strconv.ParseFloat(sensorReading.Temperature, 64)
	newAvgTemperature := ((avgTemperature * float64(totalNumOfReadings)) + temperatureReading) / float64(totalNumOfReadings+1)

	updateSensorDataWithNewReadingToRedisCache(sensorReading.SensorID,
		strconv.FormatFloat(newAvgTemperature, 'f', -1, 64), strconv.Itoa(totalNumOfReadings+1))

	sendTemperatureResponseToTheClient(w)
}
func updateSensorDataWithNewReadingToRedisCache(sensorID string, averageTemperature string, numberOfReceivedReadings string) {
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

func sendTemperatureResponseToTheClient(w http.ResponseWriter) {
	sensorAverages := make(map[string]float64)
	var total float64
	var count int

	iter := Config.RedisDb.Scan(Config.Ctx, 0, "sensor*", 0).Iterator()
	for iter.Next(Config.Ctx) {
		key := iter.Val()

		val, err := Config.RedisDb.HMGet(Config.Ctx, key, Config.AVERAGE_TEMPERATURE).Result()
		if err != nil {
			fmt.Println("Error occurred while trying to fetch field values of Sensor from redis for key: ", key)
		}

		if val[0] != nil { // make sure it's not nil
			avgStr := val[0].(string) // Redis stores as string
			avgFloat, convErr := strconv.ParseFloat(avgStr, 64)
			if convErr != nil {
				fmt.Println("Error during parsing sensor temperature: ", convErr)
				avgFloat = 0
			}

			sensorAverages[key] = avgFloat
			total += avgFloat
			count++
		}
	}

	overall := 0.0
	if count > 0 {
		overall = total / float64(count)
	}

	//Send the response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(Utils.SensorResponse{
		OverallAverage: overall,
		SensorAverages: sensorAverages,
	})
}
