package Models

type SensorEntry struct {
	SensorID                 string
	PSKHash                  string
	AverageTemperature       float64
	NumberOfReceivedReadings int
}

func (s SensorEntry) NewSensoryEntry(sensorID string, pskHash string, averageTemp float64, numOfReceivedReadings int) SensorEntry {
	var sensorEntry SensorEntry
	sensorEntry.SensorID = sensorID
	sensorEntry.PSKHash = pskHash
	sensorEntry.AverageTemperature = averageTemp
	sensorEntry.NumberOfReceivedReadings = numOfReceivedReadings

	return sensorEntry
}
