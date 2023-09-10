// Package helm.
/*
Copyright © 2023 zcubbs https://github.com/zcubbs
*/
package helm

import (
	"fmt"
	"log"
)

func chooseLogFunc(debug bool) func(string, ...interface{}) {
	if debug {
		return debugLog
	}
	return noLog
}

func debugLog(format string, v ...interface{}) {
	format = fmt.Sprintf("[debug] %s\n", format)
	err := log.Output(2, fmt.Sprintf(format, v...))
	if err != nil {
		return
	}
}

func noLog(_ string, _ ...interface{}) {}
