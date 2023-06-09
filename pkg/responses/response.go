package responses

import (
	"encoding/json"
	"net/http"

	"github.com/RafalSalwa/interview-app-srv/pkg/models"
)

// sample data struct with additional eror details
type Data struct {
	Success bool    `json:"success"`
	Message *string `json:"message"`
}

type SuccessResponse struct {
	Data Data `json:"data"`
}

type ConflictResponse struct {
	Data Data `json:"conflict"`
}

type UserResponse struct {
	Data *models.UserResponse `json:"user"`
}

func RespondInternalServerError(w http.ResponseWriter) {
	_ = NewErrorBuilder().
		SetResponseCode(http.StatusInternalServerError).
		SetReason("Internal server error").
		SetWriter(w).
		Respond()
}

func RespondNotFound(w http.ResponseWriter) {
	response := NewNotFoundResponse()
	responseBody := marshalErrorResponse(response)
	Respond(w, http.StatusNotFound, responseBody)
}

func RespondNotAuthorized(w http.ResponseWriter, msg string) {
	w.Header().Set("WWW-Authenticate", "Basic realm=Restricted")

	errorResponse := NewUnauthorizedErrorResponse(msg)
	responseBody := marshalErrorResponse(errorResponse)

	Respond(w, http.StatusUnauthorized, responseBody)
}

func RespondConflict(w http.ResponseWriter, msg string) {
	resp := NewConflictResponse(msg)
	responseBody := marshalErrorResponse(resp)
	Respond(w, http.StatusConflict, responseBody)
}

func RespondBadRequest(w http.ResponseWriter, msg string) {
	errorResponse := NewBadRequestErrorResponse(msg)

	responseBody := marshalErrorResponse(errorResponse)
	Respond(w, http.StatusBadRequest, responseBody)
}

func Respond(w http.ResponseWriter, statusCode int, responseBody []byte) {
	setHTTPHeaders(w, statusCode)
	_, _ = w.Write(responseBody)
}

func RespondOk(w http.ResponseWriter) {
	setHTTPHeaders(w, http.StatusOK)
	_, err := w.Write([]byte("{\"status\":\"ok\"}"))

	if err != nil {
		RespondInternalServerError(w)
	}
}

func NewUserResponse(u *models.UserResponse, w http.ResponseWriter) {
	response := &UserResponse{Data: u}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	js, err := json.MarshalIndent(response, "", "   ")
	if err != nil {
		RespondInternalServerError(w)
	}

	Respond(w, http.StatusOK, js)
}

func setHTTPHeaders(w http.ResponseWriter, statusCode int) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(statusCode)
}

func marshalErrorResponse(err interface{}) []byte {
	body, _ := json.Marshal(err)

	return body
}
