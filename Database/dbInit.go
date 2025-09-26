package Database

import (
	"database/sql"
	Config "dependencies/Constants"
)

func DbInit() {
	var dbErr error
	Config.Db, dbErr = sql.Open("sqlite3", "./data.Db")
	if dbErr != nil {
		panic(dbErr)
	}

	//Create Sensors table if not exists
	sensorsTableCreationQuery := `CREATE TABLE IF NOT EXISTS sensors (
    	ID INTEGER PRIMARY KEY AUTOINCREMENT,
    	SensorID Text NOT NULL UNIQUE,
		PSKHash TEXT NOT NULL
	)`

	_, dbErr = Config.Db.Exec(sensorsTableCreationQuery)
	if dbErr != nil {
		panic(dbErr)
	}

	//Create Readings table if not exists
	readingsTableCreationQuery := `CREATE TABLE IF NOT EXISTS readings
    (
		ID INTEGER PRIMARY KEY AUTOINCREMENT, 
		SensorID TEXT NOT NULL UNIQUE, 
		AverageValue FLOAT, 
		TotalNumberOfEntries INTEGER, 
		LastReadingTimestamp DATETIME DEFAULT CURRENT_TIMESTAMP, 
		FOREIGN KEY (SensorID) REFERENCES Sensors(SensorID)
    )`

	_, dbErr = Config.Db.Exec(readingsTableCreationQuery)
	if dbErr != nil {
		panic(dbErr)
	}
}
