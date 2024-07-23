package writer

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"learninggd/reader"
	"log"
	"os"
	"time"

	_ "github.com/lib/pq"
)

type Student struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Major string `json:"major"`
}

func Writer() {

	conn, ch, q, err := reader.ConnectRabbitMQ()
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer conn.Close()
	defer ch.Close()

	msgs, err := ch.Consume(
		q.Name, // queue name
		"",     // consumer tag
		true,   // auto acknowledge
		false,  // exclusive access
		false,  // no local
		false,  // no wait
		nil,    // additional arguments
	)
	if err != nil {
		log.Fatalf("Failed to register a consumer: %v", err)
	}

	db, err := ConnectPostgres()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	forever := make(chan bool)

	func() {
		for st := range msgs {
			var student Student
			err := json.Unmarshal(st.Body, &student)
			if err != nil {
				log.Printf("Error decoding JSON: %v", err)
				continue
			}

			insertStatement := `INSERT INTO Student (id, name, major)
				VALUES ($1, $2, $3)`

			_, err = db.Exec(insertStatement, student.ID, student.Name, student.Major)
			if err != nil {
				fmt.Println("Data Inserting error", err)
				return
			}

			fmt.Printf("Student Details has inserted successfully: %v, %v, %v\n", student.ID, student.Name, student.Major)

		}
	}()

	// Log message to indicate the service is waiting for messages
	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	// Block the main function to keep the service running
	<-forever
}

func ConnectPostgres() (db *sql.DB, err error) {
	maxAttempts := 5

	host := os.Getenv("DATABASE_HOST")
	port := os.Getenv("DATABASE_PORT")
	user := os.Getenv("DATABASE_USER")
	password := os.Getenv("DATABASE_PASSWORD")
	dbname := os.Getenv("DATABASE_DBNAME")
	dsn := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	for i := 0; i < maxAttempts; i++ {
		if i != 0 {
			time.Sleep(5 * time.Second)
		}
		db, err = sql.Open("postgres", dsn)
		if err != nil {
			continue
		}
		err = db.Ping()
		if err != nil {
			continue
		}
		return
	}
	return
}
