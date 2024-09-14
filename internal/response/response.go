package response

import (
	"fmt"
	"log"
	"net/http"

	"sarath/backend_project/internal/json"
)

type ResponseWriter struct{
  logger *log.Logger
}

func NewResponseWriter(logger *log.Logger) *ResponseWriter {
  return &ResponseWriter{logger}
}

func (app *ResponseWriter) logError(_ *http.Request, err error) {
  app.logger.Println(err)
}

func (app *ResponseWriter) errorResponse(w http.ResponseWriter, _ *http.Request, status int, message interface{}) {
	// getting the error msg
	errMsg := json.Envelope{"error": message}

	// writing the error msg
	err := json.WriteJSON(errMsg, w, status, nil)
	// if the message writer fails log it
	if err != nil {
		app.logError(nil, err)
		w.WriteHeader(500)
	}
}

func (app *ResponseWriter) serverErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.logError(r, err)
	message := "the server encountered a problem and could not process your request"
	app.errorResponse(w, r, http.StatusInternalServerError, message)
}

func (app *ResponseWriter) notFoundResponse(w http.ResponseWriter, _ *http.Request) {
	msg := "the requested resource couldn't be found"
	app.errorResponse(w, nil, http.StatusNotFound, msg)
}

func (app *ResponseWriter) methodNotAllowedResponse(w http.ResponseWriter, r *http.Request) {
	msg := fmt.Sprintf("the method %s isn't supported for this resource", r.Method)
	app.errorResponse(w, nil, http.StatusNotFound, msg)
}

func (app *ResponseWriter) failedValidationResponse(w http.ResponseWriter, r *http.Request, errors map[string]string) {
	app.errorResponse(w, r, http.StatusUnprocessableEntity, errors)
}
