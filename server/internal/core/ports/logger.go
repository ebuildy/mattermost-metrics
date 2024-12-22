package ports

// Logger define a simple logger, implemented by pluginAPIClient.Log
type Logger interface {
	Debug(message string, keyValuePairs ...interface{})
	Info(message string, keyValuePairs ...interface{})
	Warn(message string, keyValuePairs ...interface{})
	Error(message string, keyValuePairs ...interface{})
}
