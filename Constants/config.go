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
