package main

import (
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
)

func sendResponse(w http.ResponseWriter, message string, params ...int) {
	statusCode := http.StatusOK

	if len(params) > 0 {
		statusCode = params[0]
		w.WriteHeader(statusCode)
	}

	w.Write([]byte(message))
}

func sendServerError(w http.ResponseWriter) {
	sendResponse(w, "internal server error", http.StatusInternalServerError)
}

// LogError log an error
func LogError(err error, context ...map[string]interface{}) bool {
	if err == nil {
		return false
	}

	if len(context) > 0 {
		log.WithFields(context[0]).Error(err.Error())
	} else {
		log.Error(err.Error())
	}
	return true
}

// Prints the duration of handling the function
func printProcessingDuration(startTime time.Time) {
	dur := time.Since(startTime)

	if dur < 1500*time.Millisecond {
		log.Debugf("Duration: %s\n", dur.String())
	} else if dur > 1500*time.Millisecond {
		log.Warningf("Duration: %s\n", dur.String())
	}
}
