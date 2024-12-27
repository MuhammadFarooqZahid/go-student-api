package student

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/MuhammadFarooqZahid/go-student-api/internal/storage"
	"github.com/MuhammadFarooqZahid/go-student-api/internal/types"
	"github.com/MuhammadFarooqZahid/go-student-api/internal/utils/response"
	"github.com/go-playground/validator/v10"
)

func New(storage storage.Storage) http.HandlerFunc {
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

		lastId, err := storage.CreateStudent(
			student.Name,
			student.Email,
			student.Age,
		)

		if err != nil {
			response.WriteJson(w, http.StatusInternalServerError, err)
		}

		slog.Info("Creating a student", slog.String("id", fmt.Sprint(lastId)))
		response.WriteJson(w, http.StatusCreated, map[string]int64{"id": lastId})
	}
}

func GetById(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		slog.Info("Getting a student", slog.String("id", id))

		intId, err := strconv.ParseInt(id,10,64)

		if err != nil {
			response.WriteJson(w, http.StatusBadRequest, response.GenralError(err))
			return
		}

		student, err := storage.GetStudentById(intId)

		if err != nil {
			response.WriteJson(w, http.StatusInternalServerError, response.GenralError(err))
			return
		}

		response.WriteJson(w, http.StatusOK, student)
	}
}
