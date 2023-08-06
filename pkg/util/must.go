package util

import (
	"fmt"
	"os"
)

func Must(err error) {
	if err != nil {
		GetTheHeckOut(err)
	}
}

func GetTheHeckOut(err error) {
	fmt.Println(err)
	os.Exit(1)
}
