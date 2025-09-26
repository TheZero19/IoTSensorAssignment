package Database

type sensor struct {
	sensorID string
	pskHash  string
}

type reading struct {
	sensorID               string
	lastTemperatureReading float32
	averageTemperature     float32
	totalNumberOfEntries   int
}
