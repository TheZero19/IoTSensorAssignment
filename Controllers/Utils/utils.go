package Utils

import (
	Config "dependencies/Constants"
	"fmt"
	"net/http"
)

func CheckPOSTRequestValidity(w http.ResponseWriter, r *http.Request) bool {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return false
	}
	return true
}

func SetSensorAsDirty(sensorID string) {
	added, dirtyAddErr := Config.RedisDb.SAdd(Config.Ctx, Config.DIRTY_SENSORS_KEY, sensorID).Result()
	if added == 0 {
		fmt.Println("SensorReading already tracked for future sync with Postgres")
	} else {
		fmt.Println("SensorReading added for future sync with Postgres, Added Field Count:", added)
	}
	if dirtyAddErr != nil {
		fmt.Println("Redis SAdd Error: ", dirtyAddErr)
	}
}
