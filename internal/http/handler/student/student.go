package student

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"github.com/MuhammadFarooqZahid/go-student-api/internal/types"
	"github.com/MuhammadFarooqZahid/go-student-api/internal/utils/response"
	"github.com/go-playground/validator/v10"
)

func New() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var student types.Student

		err := json.NewDecoder(r.Body).Decode(&student)

		if errors.Is(err, io.EOF) {
			response.WriteJson(w, http.StatusBadRequest, response.GenralError(fmt.Errorf("empty body")))
			return
		}

		if err != nil {
			response.WriteJson(w, http.StatusBadRequest, response.GenralError(err))
			return
		}

		// request validation
		if err := validator.New().Struct(student); err != nil {
			validateError := err.(validator.ValidationErrors)
			response.WriteJson(w, http.StatusBadRequest, response.ValidationError(validateError))
			return
		}

		slog.Info("Creating a student")
		response.WriteJson(w, http.StatusCreated, map[string]string{"success": "Student created successfully"})
	}
}
