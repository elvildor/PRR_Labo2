package common

import (
	"log"
	"strings"
	"time"
)

const ServerAddress string = "localhost:8000"
var ServerTimeout, _ = time.ParseDuration("3s")

type Request struct {
	Header string
	Body   map[string]interface{}
}

func RemoveEndline(myString string) string {
	newString := strings.TrimSuffix(myString, "\n")
	newString = strings.TrimSuffix(newString, "\r")
	return newString
}

func Debug(message string) {
	log.Println("[DEBUG] " + time.Now().String() + " : " + message)
}

func Trace(message string) {
	log.Println("[TRACE] " + time.Now().String() + " : " + message)
}


