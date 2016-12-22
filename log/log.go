// Scratche Backend - Logger.
//
// Logger helpers.
//
// Rémy Mathieu © 2016

package log

import (
	"log"
)

func Error(data ...interface{}) {
	log.Println("Error", data)
}

func Warning(data ...interface{}) {
	log.Println("Warning", data)
}

func Info(data ...interface{}) {
	log.Println("Info", data)
}

func Debug(data ...interface{}) {
	log.Println("Debug", data)
}
