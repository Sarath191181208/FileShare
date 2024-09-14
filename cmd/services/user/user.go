package user

import (
	"errors"
	"log"
	"net/http"

	"sarath/backend_project/internal/data"
	"sarath/backend_project/internal/json"
	"sarath/backend_project/internal/response"
	"sarath/backend_project/internal/validator"
)

type Handler struct {
	Logger *log.Logger
	models *data.Models
}

func NewHandler(logger *log.Logger, models *data.Models) *Handler {
	return &Handler{logger, models}
}

func (h *Handler) RegisterUserHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	responseWriter := response.NewResponseWriter(h.Logger)

	// read the input data
	err := json.ReadJSON(&input, w, r)
	if err != nil {
		responseWriter.BadRequestResponse(w, r, err)
		return
	}

	// create a new user with the input data
	user := &data.User{
		Email: input.Email,
	}
  // try hash the password
	err = user.Password.Set(input.Password)
	if err != nil {
		responseWriter.ServerErrorResponse(w, r, err)
		return
	}

	// validate the user data
	v := validator.New()
	if data.ValidateUser(v, user); !v.Valid() {
		responseWriter.FailedValidationResponse(w, r, v.Errors)
		return
	}

	// try inserting the user into the database
	err = h.models.Users.Insert(user)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrDuplicateEmail):
			v.AddError("email", "a user with this email address already exists")
			responseWriter.FailedValidationResponse(w, r, v.Errors)
		default:
			responseWriter.ServerErrorResponse(w, r, err)
		}
		return
	}

	// send the user data in the response
	err = json.WriteJSON(json.Envelope{"user": user}, w, http.StatusCreated, nil)
	if err != nil {
		responseWriter.ServerErrorResponse(w, r, err)
	}
}

func (h *Handler) LoginUserHandler(w http.ResponseWriter, r *http.Request) {
	//
}
