package Database

import (
	"database/sql"
	Config "dependencies/Constants"
	"fmt"
	"log"

	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
)

func DbInit() {
	isPostgresInit := postgresInit()
	isRedisInit := redisInit()

	if isPostgresInit {
		onPostgresInit()
	}

	if isRedisInit {
		onRedisInit()
	}
}

func postgresInit() bool {
	var dbErr error

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		Config.Env.PostgresUser,
		Config.Env.PostgresPassword,
		Config.Env.PostgresHost,
		Config.Env.PostgresPort,
		Config.Env.PostgresDB)

	Config.PostgresDb, dbErr = sql.Open("postgres", dsn)
	if dbErr != nil {
		panic(dbErr)
		log.Println("Failed to connect to database: %v", dbErr)
		Config.PostgresDb.Close()
		return false
	}

	return true
}
func onPostgresInit() {
	//Create Sensors table if not exists
	sensorsTableCreationQuery := `CREATE TABLE IF NOT EXISTS sensors (
    	ID SERIAL PRIMARY KEY,
    	sensor_id Text NOT NULL UNIQUE,
		psk_hash TEXT NOT NULL,
		average_temperature FLOAT NOT NULL DEFAULT 0.0,
		num_of_received_readings INTEGER NOT NULL DEFAULT 0
	)`

	_, dbErr := Config.PostgresDb.Exec(sensorsTableCreationQuery)
	if dbErr != nil {
		fmt.Println("Error creating sensors table: %v", dbErr)
	}
}

func redisInit() bool {
	redisAddr := fmt.Sprintf("%s:%s", Config.Env.RedisHost, Config.Env.RedisPort)
	Config.RedisDb = redis.NewClient(&redis.Options{
		Addr: redisAddr,
	})
	if _, err := Config.RedisDb.Ping(Config.Ctx).Result(); err != nil {
		panic(err)
		return false
	}

	fmt.Println("Successfully connected to Redis")
	return true
}

func onRedisInit() {
	//Add instructions here to execute when redis is initialized
}
