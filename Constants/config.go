package Constants

import (
	"context"
	"database/sql"

	"github.com/redis/go-redis/v9"
)

var PostgresDb *sql.DB
var RedisDb *redis.Client
var Ctx = context.Background()

const AUTH_HEADER_KEY string = "Authorization"
const AUTH_HEADER_VALUE_SEPARATOR string = " "
const API_KEY_AUTH_REGISTRATION_TYPE_PREFIX string = "API-KEY"
const PSK_AUTH_TYPE_PREFIX string = "PSK"

// Redis Keys for Hash and Sets
const DIRTY_SENSORS_KEY = "dirty_sensors"
const PSK_HASH = "PSKHash"
const AVERAGE_TEMPERATURE = "AverageTemperature"
const NUMBER_OF_RECEIVED_READINGS = "NumberOfReceivedReadings"

// Environment Variables Names
const POSTGRES_USER = "POSTGRES_USER"
const POSTGRES_PASSWORD = "POSTGRES_PASSWORD"
const POSTGRES_DB = "POSTGRES_DB"
const POSTGRES_HOST = "POSTGRES_HOST"
const POSTGRES_PORT = "POSTGRES_PORT"
const REDIS_HOST = "REDIS_HOST"
const REDIS_PORT = "REDIS_PORT"
const SENSOR_REGISTRATION_API_KEY = "SENSOR_REGISTRATION_API_KEY"

const SYNC_FACTOR = "SYNC_FACTOR"

// Environment variables
type EnvironmentVariables struct {
	PostgresUser             string
	PostgresPassword         string
	PostgresDB               string
	PostgresHost             string
	PostgresPort             string
	RedisHost                string
	RedisPort                string
	DatabaseSyncFactor       int
	SensorRegistrationAPIKey string
}

var Env EnvironmentVariables
