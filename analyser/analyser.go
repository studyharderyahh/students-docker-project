package analyser

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

type Student struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Major string `json:"major"`
}


func fetchStudents(apiURL string) ([]Student, error) {
	// Call the API to get the students
	resp, err := http.Get(apiURL)
	if err != nil {
		return nil, fmt.Errorf("failed to call API: %v", err)
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read API response: %v", err)
	}

	// Unmarshal the response body into a slice of students
	students := []Student{}
	err = json.Unmarshal(body, &students)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal API response: %v", err)
	}

	return students, nil
}

func writeStudentsToFile(students []Student, filePath string) error {
	// Create and write to the file
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create file: %v", err)
	}
	defer file.Close()

	for _, student := range students {
		fmt.Printf("Student Details read successfully: %v, %v, %v\n", student.ID, student.Name, student.Major)
		studentID := strconv.Itoa(student.ID)
		studentData := []string{studentID, student.Name, student.Major}
		_, err := file.WriteString(strings.Join(studentData, ", ") + "\n")
		if err != nil {
			return fmt.Errorf("failed to write to file: %v", err)
		}
	}

	return nil
}

func Analyser() {
	apiHost := os.Getenv("API_HOST")
	apiPort := os.Getenv("API_PORT")
	apiEndPoint := os.Getenv("API_ENDPOINT")
	apiURL := fmt.Sprintf("http://%s:%s/%s", apiHost, apiPort, apiEndPoint)

	students, err := fetchStudents(apiURL)
	if err != nil {
		log.Fatalf("Error fetching students: %v", err)
	}

	filePath := os.Getenv("FILE_PATH")
	err = writeStudentsToFile(students, filePath)
	if err != nil {
		log.Fatalf("Error writing students to file: %v", err)
	}

	fmt.Printf("Data successfully written to file: %s\n", filePath)
}