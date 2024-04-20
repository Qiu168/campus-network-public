package main

import (
	"os"
	"strings"
)

func createFileIfNotExistAndWrite(filePath, usernameText, passwordText string) error {
	file, err := os.OpenFile(filePath, os.O_CREATE, 0777)
	defer func() { file.Close() }()
	if err == nil {
		_, err = file.WriteString(usernameText + "\n" + passwordText)
	}
	return err
}

func getConfig() (string, string, error) {
	file, err := os.ReadFile("./config.txt")
	if err != nil {
		return "", "", err
	}
	config := string(file)
	split := strings.Split(config, "\n")
	return split[0], split[1], nil
}
