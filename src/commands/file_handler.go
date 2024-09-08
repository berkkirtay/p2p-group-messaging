// Copyright (c) 2024 Berk Kirtay

package commands

import (
	"os"
)

func dumpToFile(data string, fileName string) {
	file, err := os.Create(fileName)
	if err != nil {
		panic(err)
	}
	_, err = file.Write([]byte(data))
	if err != nil {
		panic(err)
	}
	defer file.Close()
}

func readFromFile(fileName string) string {
	data, err := os.ReadFile(fileName)
	if err != nil {
		panic(err)
	}
	return string(data)
}
