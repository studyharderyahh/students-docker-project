package students

import (
	"database/sql"
	"encoding/json"
	"learninggd/writer"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"
)

type Student struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Major string `json:"major"`
}

func StudentApi() {
	e := echo.New()
	e.GET("/health", func(c echo.Context) error {
		return c.NoContent(http.StatusOK)
	})

	e.GET("/students", getAllStudents)
	e.GET("/student", handleGetStudent)
	e.POST("/students", handleCreateStudent)
	e.PUT("/students", handleUpdateStudent)
	e.DELETE("/students", deleteStudent)

	// Start the server
	e.Logger.Fatal(e.Start(":8080"))
}

func handleGetStudent(c echo.Context) error {
	idParam := c.QueryParam("id")
	nameParam := c.QueryParam("name")
	majorParam := c.QueryParam("major")

	if idParam != "" {
		id, err := strconv.Atoi(idParam)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid ID"})
		}
		return getStudentByID(c, id)
	} else if nameParam != "" {
		return getStudentByName(c, nameParam)
	} else if majorParam != "" {
		return getStudentByMajor(c, majorParam)
	}

	// If no query parameter is provided, return all students
	return getAllStudents(c)
}

func getAllStudents(c echo.Context) error {
	db, err := writer.ConnectPostgres()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	defer db.Close()

	rows, err := db.Query("SELECT id, name, major FROM Student")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	defer rows.Close()

	students := []Student{}
	for rows.Next() {
		var st Student
		if err := rows.Scan(&st.ID, &st.Name, &st.Major); err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}
		students = append(students, st)
	}

	return c.JSON(http.StatusOK, students)
}

func getStudentByID(c echo.Context, studentID int) error {
	db, err := writer.ConnectPostgres()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	defer db.Close()

	var st Student
	err = db.QueryRow("SELECT id, name, major FROM Student WHERE id = $1", studentID).Scan(&st.ID, &st.Name, &st.Major)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.JSON(http.StatusNotFound, map[string]string{"message": "Student not found"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, st)
}

func getStudentByName(c echo.Context, studentName string) error {
	db, err := writer.ConnectPostgres()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	defer db.Close()

	rows, err := db.Query("SELECT id, name, major FROM Student WHERE name = $1", studentName)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	defer rows.Close()

	students := []Student{}
	for rows.Next() {
		var st Student
		if err := rows.Scan(&st.ID, &st.Name, &st.Major); err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}
		students = append(students, st)
	}

	if len(students) == 0 {
		return c.JSON(http.StatusNotFound, map[string]string{"message": "No students found with that name"})
	}

	return c.JSON(http.StatusOK, students)
}

func getStudentByMajor(c echo.Context, studentMajor string) error {
	db, err := writer.ConnectPostgres()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	defer db.Close()

	rows, err := db.Query("SELECT id, name, major FROM Student WHERE major = $1", studentMajor)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	defer rows.Close()

	students := []Student{}
	for rows.Next() {
		var st Student
		if err := rows.Scan(&st.ID, &st.Name, &st.Major); err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}
		students = append(students, st)
	}

	if len(students) == 0 {
		return c.JSON(http.StatusNotFound, map[string]string{"message": "No students found with that major"})
	}

	return c.JSON(http.StatusOK, students)
}

func handleCreateStudent(c echo.Context) error {
	var st Student
	if err := parseRequestBody(c, &st); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	if st.ID == 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Student ID is required"})
	}

	return createStudent(c, st)
}

func createStudent(c echo.Context, st Student) error {
	db, err := writer.ConnectPostgres()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	defer db.Close()

	_, err = db.Exec("INSERT INTO Student (id, name, major) VALUES ($1, $2, $3)", st.ID, st.Name, st.Major)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, st)
}

func handleUpdateStudent(c echo.Context) error {
	var st Student
	if err := parseRequestBody(c, &st); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	if st.ID == 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Student ID is required"})
	}

	return updateStudent(c, st)
}

func updateStudent(c echo.Context, st Student) error {
	db, err := writer.ConnectPostgres()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	defer db.Close()

	result, err := db.Exec("UPDATE Student SET name = $2, major = $3 WHERE id = $1", st.ID, st.Name, st.Major)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	if rowsAffected == 0 {
		return c.JSON(http.StatusNotFound, map[string]string{"message": "Student not found"})
	}

	return c.JSON(http.StatusOK, st)
}

func deleteStudent(c echo.Context) error {
	var st Student
	if err := parseRequestBody(c, &st); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	if st.ID == 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Student ID is required"})
	}

	db, err := writer.ConnectPostgres()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	defer db.Close()

	result, err := db.Exec("DELETE FROM Student WHERE id = $1", st.ID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	if rowsAffected == 0 {
		return c.JSON(http.StatusNotFound, map[string]string{"message": "Student not found"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Student deleted successfully"})
}

// parseRequestBody reads and parses the request body
func parseRequestBody(c echo.Context, body interface{}) error {
	if err := json.NewDecoder(c.Request().Body).Decode(body); err != nil {
		return err
	}
	return nil
}
