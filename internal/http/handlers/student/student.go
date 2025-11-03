package student

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/manish-npx/go-student-api/internal/types"
	"github.com/manish-npx/go-student-api/internal/utils/response"
)

func New() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var student types.Student

		err := json.NewDecoder(r.Body).Decode(&student)
		if errors.Is(err, io.EOF) {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(fmt.Errorf("empty body")))
			return
		}
		if err != nil {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(fmt.Errorf("invalid JSON: %v", err)))
			return
		}

		//request validation

		if err := validator.New().Struct(student); err != nil {
			//type casting err.(validator.ValidationErrors)
			response.WriteJson(w, http.StatusBadRequest, response.ValidationError(err.(validator.ValidationErrors)))
			return
		}

		data := map[string]any{
			"success": true,
			"student": student,
			"message": "Student record created successfully",
		}

		slog.Info("Creating student record",
			slog.String("name", student.Name),
			slog.String("email", student.Email),
		)

		response.WriteJson(w, http.StatusCreated, data)

	}

}
