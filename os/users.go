package os

import (
	"bufio"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"github.com/pkg/errors"
	"os"
	"os/exec"
	"strings"
)

const (
	userFile string = "/etc/passwd"
)

type User struct {
	Name      string `json:"name"`
	Directory string `json:"directory"`
	Group     string `json:group`
	Shell     string `json:shell`
}

// ReadEtcPasswd file /etc/passwd and return slice of users
func ReadEtcPasswd() (list []string) {

	file, err := os.Open(userFile)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	r := bufio.NewScanner(file)

	for r.Scan() {
		lines := r.Text()
		parts := strings.Split(lines, ":")
		list = append(list, parts[0])
	}
	return list
}

// check if user on the host
func check(s []string, u string) bool {
	for _, w := range s {
		if u == w {
			return true
		}
	}
	return false
}

// Return securely generated random bytes

func CreateRandom(n int) string {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		fmt.Println(err)
		//os.Exit(1)
	}
	return string(b)
}

// AddUserIfNotExist check if user exist on the host
func AddUserIfNotExist(u *User) (string, error) {

	users := ReadEtcPasswd()

	if check(users, u.Name) {
		return "", errors.New("User already exists")
	} else {
		return AddNewUser(u)
	}
}

// AddNewUser is created by executing shell command useradd
func AddNewUser(u *User) (string, error) {

	encrypt := base64.StdEncoding.EncodeToString([]byte(CreateRandom(9)))

	argUser := []string{"--disabled-password", "--gecos", "\"\"", u.Name}
	argPass := []string{"-c", fmt.Sprintf("echo %s:%s | chpasswd", u.Name, encrypt)}

	userCmd := exec.Command("adduser", argUser...)
	passCmd := exec.Command("/bin/sh", argPass...)

	if out, err := userCmd.Output(); err != nil {
		return "", err
	} else {

		fmt.Printf("Output: %s\n", out)

		if _, err := passCmd.Output(); err != nil {
			fmt.Println(err)
			return "", err
		}
		return encrypt, nil
	}
}

// DeleteUserIfExist check if user exist on the host
func DeleteUserIfExist(name string) error {
	users := ReadEtcPasswd()

	if check(users, name) {
		return DeleteUser(name)
	} else {
		return errors.New("user doesn't exists")
	}
}

// DeleteUser is created by executing shell command userdel
func DeleteUser(name string) error {
	arg := []string{name}

	cmd := exec.Command("deluser", arg...)

	if out, err := cmd.Output(); err != nil {
		return err
	} else {
		fmt.Printf("Output: %s\n", out)
		return nil
	}
}
