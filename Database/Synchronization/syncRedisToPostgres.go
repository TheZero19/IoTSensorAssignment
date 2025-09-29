package Synchronization

import (
	Config "dependencies/Constants"
	"fmt"
	"log"
	"strconv"
	"time"
)

func syncDirtySensors() {
	sensorIDs, err := Config.RedisDb.SMembers(Config.Ctx, Config.DIRTY_SENSORS_KEY).Result()

	if err != nil {
		fmt.Println("Error while fetching the dirty sensors: ", err)
		return
	}

	if len(sensorIDs) == 0 {
		fmt.Println("No sensors to be synchronized currently...")
		return
	}

	fmt.Println("Sensors to be synchronized. Count: ", len(sensorIDs))

	for _, sensorID := range sensorIDs {
		vals, err := Config.RedisDb.HMGet(
			Config.Ctx,
			sensorID,
			Config.PSK_HASH,
			Config.AVERAGE_TEMPERATURE,
			Config.NUMBER_OF_RECEIVED_READINGS).Result()

		pskHash := vals[0].(string)
		averageTemperature, floatParseErr := strconv.ParseFloat((vals[1]).(string), 64)
		if floatParseErr != nil {
			fmt.Println("Average temperature parse error: ", floatParseErr)
			continue
		}

		totalNumOfReadings, intParseErr := strconv.Atoi(vals[2].(string))
		if intParseErr != nil {
			fmt.Println("NumberOfReadings parse error: ", intParseErr)
			continue
		}

		if err != nil {
			log.Println("Error while fetching the data for sensorID: ", sensorID, err)
			continue
		}

		postgresQuery := `INSERT INTO sensors 
							(sensor_id, 
							psk_hash, 
							average_temperature,
							num_of_received_readings) 
							VALUES ($1, $2, $3, $4)
							ON CONFLICT(sensor_id) DO UPDATE SET
							average_temperature=EXCLUDED.average_temperature,
							num_of_received_readings=EXCLUDED.num_of_received_readings`
		_, err = Config.PostgresDb.Exec(postgresQuery, sensorID, pskHash, averageTemperature, totalNumOfReadings)

		if err != nil {
			log.Println("Error while syncing sensor with ID: ", sensorID, err)
			continue
		}

		Config.RedisDb.SRem(Config.Ctx, Config.DIRTY_SENSORS_KEY, sensorID)
	}
}

func StartBackgroundSync(interval time.Duration) {
	fmt.Println("Starting background sync")
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for {
			<-ticker.C
			syncDirtySensors()
		}
	}()
}
