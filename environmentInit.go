package main

import (
	Config "dependencies/Constants"
	"fmt"
	"os"
	"strconv"
)

func environmentVarInit() {
	Config.Env.PostgresUser = os.Getenv(Config.POSTGRES_USER)
	Config.Env.PostgresPassword = os.Getenv(Config.POSTGRES_PASSWORD)
	Config.Env.PostgresDB = os.Getenv(Config.POSTGRES_DB)
	Config.Env.PostgresHost = os.Getenv(Config.POSTGRES_HOST)
	Config.Env.PostgresPort = os.Getenv(Config.POSTGRES_PORT)
	Config.Env.RedisHost = os.Getenv(Config.REDIS_HOST)
	Config.Env.RedisPort = os.Getenv(Config.REDIS_PORT)
	Config.Env.SensorRegistrationAPIKey = os.Getenv(Config.SENSOR_REGISTRATION_API_KEY)

	//For sync factor
	syncFactor := os.Getenv(Config.SYNC_FACTOR)
	factor, err := strconv.Atoi(syncFactor)
	if err != nil {
		fmt.Println("Error converting sync factor to int", err)
		factor = 10 //Setting default factor of 10 seconds for database sync
	}
	Config.Env.DatabaseSyncFactor = factor

	fmt.Println("Environment variables initialized, with database sync factor (in seconds): ", Config.Env.DatabaseSyncFactor)
}
