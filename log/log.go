// Scratche Backend - Logger.
//
// Logger helpers.
//
// Rémy Mathieu © 2016

package log

import (
	"fmt"
	"log"

	"github.com/fatih/color"
)

var (
	yellow = color.New(color.FgYellow).SprintFunc()
	red    = color.New(color.FgRed).SprintFunc()
	blue   = color.New(color.FgBlue).SprintFunc()
	green  = color.New(color.FgGreen).SprintFunc()
)

func Error(data ...interface{}) {
	log.Println(red("Error"), data)
}

func Warning(data ...interface{}) {
	log.Println(yellow("Warning"), data)
}

func Info(data ...interface{}) {
	log.Println(green("Info"), data)
}

func Debug(data ...interface{}) {
	log.Println(blue("Debug"), data)
}

func Err(prefix string, err error) error {
	return fmt.Errorf("%s %v", err)
}
