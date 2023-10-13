package main

import (
	"os"
	"strings"
)

func getConfig() (string, string) {
	file, _ := os.ReadFile("./config.txt")
	config := string(file)
	split := strings.Split(config, "\n")
	return split[0][:len(split[0])-1], split[1]
}
