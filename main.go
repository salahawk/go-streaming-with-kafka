package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/joho/godotenv"
)

type Message struct {
	Name string  `json:"name"`
	Size float64 `json:"size"`
	Time int64   `json:"time"`
}

func iterate(path string) error {
	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Fatalf(err.Error())
			return err
		}
		m := Message{
			Name: info.Name(),
			Size: float64(info.Size() / (1 << 20)),
			Time: time.Now().Unix(),
		}

		v, err := json.Marshal(m)
		if err != nil {
			return err
		}
		fmt.Printf("File Name: %v\n", string(v))
		return nil
	})
	return err
}

func main() {
	godotenv.Load()

	currentDirectory, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	for {
		err := iterate(currentDirectory + "/files")
		if err != nil {
			fmt.Printf("Failed to run sample: %v\n", err)
		}
		time.Sleep(time.Minute)
	}
}
