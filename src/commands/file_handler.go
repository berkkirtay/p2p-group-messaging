// Copyright (c) 2024 Berk Kirtay

package commands

import (
	"os"
	"path/filepath"
)

func dumpToFile(data string, fileName string) {
	dir := filepath.Dir(fileName)
	err := os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		panic(err)
	}
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
