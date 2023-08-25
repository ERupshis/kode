package logger

import "net/http"

type BaseLogger interface {
	Info(msg string, fields ...interface{})
	LogHandler(h http.Handler) http.Handler
	Sync()
}

func CreateZapLogger(level string) BaseLogger {
	log, err := createZapLogger(level)
	if err != nil {
		panic(err)
	}

	return log
}
