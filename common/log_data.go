package common

func NewLogData(collectEntryList []CollectEntry) (logData *LogData) {
	logData = &LogData{}
	logData.Lock.Lock()
	logData.CollectEntryList = collectEntryList
	logData.Lock.Unlock()
	return logData
}
