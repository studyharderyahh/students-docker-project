package reader

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	_ "github.com/lib/pq"
	amqp "github.com/rabbitmq/amqp091-go"
)

type Student struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Major string `json:"major"`
}

var students = []Student{}

var myfilePath string = "./TestFile/studentInfoproducer.txt"

func Reader() {
	readFromFile()

	conn, ch, q, err := ConnectRabbitMQ()
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer conn.Close()
	defer ch.Close()

	for _, st := range students {
		body, err := json.Marshal(st)
		if err != nil {
			log.Fatalf("Failed to marshal JSON: %v", err)
		}
		err = ch.Publish(
			"",
			q.Name,
			false,
			false,
			amqp.Publishing{
				ContentType: "application/json",
				Body:        body, // No need to convert body to []byte again
			})
		if err != nil {
			log.Fatalf("Failed to publish message: %v", err)
		}
	}
}

func ConnectRabbitMQ() (conn *amqp.Connection, ch *amqp.Channel, q amqp.Queue, err error) {
	rabbitMQHost := os.Getenv("RABBITMQ_HOST")
	rabbitMQPort := os.Getenv("RABBITMQ_PORT")
	rabbitMQUser := os.Getenv("RABBITMQ_USER")
	rabbitMQPassword := os.Getenv("RABBITMQ_PASSWORD")
	rabbitMQURL := fmt.Sprintf("amqp://%s:%s@%s:%s/", rabbitMQUser, rabbitMQPassword, rabbitMQHost, rabbitMQPort)

	maxAttempts := 5
	for i := 0; i < maxAttempts; i++ {
		if i != 0 {
			time.Sleep(5 * time.Second)
		}
		conn, err = amqp.Dial(rabbitMQURL)
		if err != nil {
			if i == maxAttempts-1 {
				log.Fatalf("Failed to connect to RabbitMQ after multiple attempts: %v", err)
			}
			continue
		}

		ch, err = conn.Channel()
		if err != nil {
			conn.Close()
			if i == maxAttempts-1 {
				log.Fatalf("Failed to open a channel after multiple attempts: %v", err)
			}
			continue
		}

		q, err = ch.QueueDeclare(
			"student_messages", // name
			false,              // durable
			false,              // delete when unused
			false,              // exclusive
			false,              // no-wait
			nil,                // arguments
		)
		if err != nil {
			ch.Close()
			conn.Close()
			if i == maxAttempts-1 {
				log.Fatalf("Failed to declare a queue after multiple attempts: %v", err)
			}
			continue
		}
		return
	}
	return
}

func readFromFile() {
	file, err := os.Open(myfilePath)
	if err != nil {
		log.Fatalf("File reading error: %v", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		item := strings.Split(line, ",")

		if len(item) != 3 {
			log.Printf("Invalid line format: %s", line)
			continue
		}

		id, err := strconv.Atoi(item[0])
		if err != nil {
			log.Printf("Failed to convert ID to integer: %v", err)
			continue
		}

		student := Student{ID: id, Name: item[1], Major: item[2]}
		students = append(students, student)
	}

	if err := scanner.Err(); err != nil {
		log.Fatalf("Error reading from file: %v", err)
	}
}
