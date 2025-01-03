package api_utils

import (
	"encoding/json"
	"fmt"
	"net/http"

	database_models "github.com/KonovalovIly/anki_pdf/database/model"
	"github.com/go-playground/validator/v10"
)

var Validator *validator.Validate

func init() {
	Validator = validator.New(validator.WithRequiredStructEnabled())
}

func writeJson(w http.ResponseWriter, status int, data any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(data)
}

func ReadJson(w http.ResponseWriter, r *http.Request, data any) error {
	maxBytes := 1_048_578 // 1 MB limit
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	return decoder.Decode(data)
}

func WriteJsonError(w http.ResponseWriter, status int, err error) error {
	type envelope struct {
		Error string `json:"error"`
	}
	fmt.Print(err)

	return writeJson(w, status, &envelope{Error: err.Error()})
}

func WriteJsonDatabaseError(w http.ResponseWriter, status int, err *database_models.DatabaseError) error {
	type envelope struct {
		Error string `json:"error"`
		Type  string `json:"type"`
	}
	fmt.Print(err)

	return writeJson(w, status, &envelope{Error: err.Error, Type: err.Typ})
}

func JsonResponse(w http.ResponseWriter, status int, data any) error {
	type envelope struct {
		Data any `json:"data"`
	}
	return writeJson(w, status, &envelope{Data: data})
}
