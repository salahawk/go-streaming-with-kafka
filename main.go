package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/joho/godotenv"
)

type Message struct {
	Name string  `json:"name"`
	Size float64 `json:"size"`
	Time int64   `json:"time"`
}

func iterate(path string) error {
	var bootstrapServers = os.Getenv("bootstrapServers")
	var ccloudAPIKey = os.Getenv("ccloudAPIKey")
	var ccloudAPISecret = os.Getenv("ccloudAPISecret")

	fmt.Println(bootstrapServers)
	fmt.Println(ccloudAPIKey)
	fmt.Println(ccloudAPISecret)

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

		messageValue, err := json.Marshal(m)
		if err != nil {
			return err
		}
		fmt.Printf("File Name: %v\n", string(messageValue))

		// create producer & produce message
		topic := "files"
		// Produce a new record to the topic...
		producer, err := kafka.NewProducer(&kafka.ConfigMap{
			"bootstrap.servers": bootstrapServers,
			"sasl.mechanisms":   "PLAIN",
			"security.protocol": "SASL_SSL",
			"sasl.username":     ccloudAPIKey,
			"sasl.password":     ccloudAPISecret,
		})
		if err != nil {
			panic(fmt.Sprintf("Failed to create producer: %s", err))
		}

		producer.Produce(&kafka.Message{
			TopicPartition: kafka.TopicPartition{
				Topic:     &topic,
				Partition: kafka.PartitionAny,
			},
			Value: messageValue,
		}, nil)
		// Wait for delivery report
		e := <-producer.Events()

		fmt.Println(e)

		message := e.(*kafka.Message)
		if message.TopicPartition.Error != nil {
			fmt.Printf("failed to deliver message: %v\n",
				message.TopicPartition)
		} else {
			fmt.Printf("delivered to topic %s [%d] at offset %v\n",
				*message.TopicPartition.Topic,
				message.TopicPartition.Partition,
				message.TopicPartition.Offset)
		}

		producer.Close()

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
