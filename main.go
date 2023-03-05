/*
Copyright © 2023 zcubbs https://github.com/zcubbs
*/
package main

import (
	"github.com/zcubbs/zrun/cmd"
	"github.com/zcubbs/zrun/configs"
)

func init() {
	configs.Bootstrap()
}

func main() {
	cmd.Execute()
}
