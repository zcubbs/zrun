package util

import (
	"fmt"
	"os"
)

func Must(err error) {
	if err != nil {
		getTheHeckOut(err)
	}
}

func getTheHeckOut(err error) {
	fmt.Println(err)
	os.Exit(1)
}
