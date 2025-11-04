package student

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/manish-npx/go-student-api/internal/storage"
	"github.com/manish-npx/go-student-api/internal/types"
	"github.com/manish-npx/go-student-api/internal/utils/response"
)

// 🧩 POST /api/student
// ---------------------------------------------------------
// This handler creates a new student record.
// 1. Validates HTTP method (must be POST)
// 2. Decodes JSON body → types.Student
// 3. Validates fields using go-playground/validator
// 4. Calls `storage.CreateStudent()` to persist the record
// 5. Responds with JSON containing success info
func New(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// ✅ Ensure correct HTTP method
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var student types.Student

		// 🧠 Decode request body JSON → Go struct
		err := json.NewDecoder(r.Body).Decode(&student)
		if errors.Is(err, io.EOF) {
			// Empty body — client sent no JSON
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(fmt.Errorf("empty body")))
			return
		}
		if err != nil {
			// Invalid JSON syntax
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(fmt.Errorf("invalid JSON: %v", err)))
			return
		}

		// 🧩 Request validation
		// Uses struct tags in `types.Student` (e.g., validate:"required")
		if err := validator.New().Struct(student); err != nil {
			response.WriteJson(w, http.StatusBadRequest, response.ValidationError(err.(validator.ValidationErrors)))
			return
		}

		// 💾 Insert student into DB via storage layer
		lastId, err := storage.CreateStudent(
			student.Name,
			student.Email,
			student.Age,
		)
		if err != nil {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(err))
			return
		}

		// 📦 Build success response payload
		data := map[string]any{
			"success": true,
			"id":      lastId,
			"student": student,
			"message": "Student record created successfully",
		}

		// 🪵 Log structured info about the new record
		slog.Info("Creating student record",
			slog.String("name", student.Name),
			slog.String("email", student.Email),
			slog.Int64("id", lastId),
		)

		// 🚀 Send response
		response.WriteJson(w, http.StatusCreated, data)
	}
}

// 🧩 GET /api/student/{id}
// ---------------------------------------------------------
// Fetches a single student record by ID.
// 1. Extracts `id` path param
// 2. Converts string → int64
// 3. Calls `storage.GetStudentById()`
// 4. Returns the record in JSON
func GetById(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		slog.Info("Getting a student record", slog.String("id", id))

		// 🔢 Convert id from string → int64
		intId64, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(fmt.Errorf("invalid id %v", id)))
			return
		}

		// 💾 Fetch record from DB
		student, err := storage.GetStudentById(intId64)
		if err != nil {
			slog.Error("Error getting student record", slog.String("error", err.Error()))
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(err))
			return
		}

		// 🚀 Respond with found record
		response.WriteJson(w, http.StatusOK, student)
	}
}

// 🧩 GET /api/students
// ---------------------------------------------------------
// Fetches all student records.
// 1. Calls `storage.GetStudents()`
// 2. Returns array of students as JSON
func GetList(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slog.Info("Getting all student records")

		// 💾 Retrieve all students from DB
		students, err := storage.GetStudents()
		if err != nil {
			slog.Error("Error getting students", slog.String("error", err.Error()))
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(err))
			return
		}

		// 🚀 Send JSON list
		response.WriteJson(w, http.StatusOK, students)
	}
}

// 🧩 PUT /api/student/{id}
// ---------------------------------------------------------
// This handler update creates a new student record.
// 1. Validates HTTP method (must be PUT)
// 2. Decodes JSON body → types.Student
// 3. Validates fields using go-playground/validator

func UpdateById(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slog.Info("Update student record based on Id")

		// ✅ Ensure correct HTTP method
		if r.Method != http.MethodPut {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var student types.Student

		// 🧠 Decode request body JSON → Go struct
		err := json.NewDecoder(r.Body).Decode(&student)
		if errors.Is(err, io.EOF) {
			// Empty body — client sent no JSON
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(fmt.Errorf("empty body")))
			return
		}
		if err != nil {
			// Invalid JSON syntax
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(fmt.Errorf("invalid JSON: %v", err)))
			return
		}
		id := r.PathValue("id")
		slog.Info("Getting a student record", slog.String("id", id))

		// 🔢 Convert id from string → int64
		intId64, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(fmt.Errorf("invalid id %v", id)))
			return
		}

		// 💾 Retrieve all students from DB
		lastId, err := storage.UpdateStudentById(
			intId64,
			student.Name,
			student.Email,
			student.Age,
		)
		if err != nil {
			slog.Error("Error getting students", slog.String("error", err.Error()))
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(err))
			return
		}

		// 🚀 Send JSON list// 📦 Build success response payload
		data := map[string]any{
			"success": true,
			"id":      lastId.ID,
			"student": student,
			"message": "Student record created successfully",
		}

		// 🪵 Log structured info about the new record
		slog.Info("Updated student record",
			slog.String("name", student.Name),
			slog.String("email", student.Email),
		)

		// 🚀 Send response
		response.WriteJson(w, http.StatusOK, data)
	}
}
