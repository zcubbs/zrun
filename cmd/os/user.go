// Package os
/*
Copyright © 2023 zcubbs https://github.com/zcubbs
*/
package os

import (
	"fmt"
	"github.com/spf13/cobra"
	zos "github.com/zcubbs/zrun/os"
)

// addUser represents the list command
var addUser = &cobra.Command{
	Use:   "adduser",
	Short: "Add user to the system, if user doesn't exists",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Add user", args)
		u := &zos.User{
			Name:      args[0],
			Directory: "",
			Group:     "",
			Shell:     "",
		}

		passwd, err := zos.AddUserIfNotExist(u)
		if err != nil {
			panic(err)
			return
		}

		fmt.Printf("User added successfully. Password: %s\n", passwd)
	},
}

// deleteUser represents the list command
var deleteUser = &cobra.Command{
	Use:   "deluser",
	Short: "Delete user, if user exists",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Delete user", args)

		err := zos.DeleteUserIfExist(args[0])
		if err != nil {
			fmt.Println("User doesn't exists")
			return
		}

		fmt.Println("User deleted successfully")
	},
}

func init() {
	Cmd.AddCommand(addUser)
	Cmd.AddCommand(deleteUser)
}
