package user

import (
	"errors"
	"log"
	"net/http"

	"sarath/backend_project/internal/data"
	"sarath/backend_project/internal/json"
	"sarath/backend_project/internal/jwt"
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

func (h *Handler) GetLoginUserHandler(jwtSecret string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
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

    // validate the incomming data 
    validator := validator.New()

    data.ValidateEmail(validator, input.Email)
    data.ValidatePasswordPlaintext(validator, input.Password)

    if !validator.Valid() {
      responseWriter.FailedValidationResponse(w, r, validator.Errors)
      return 
    }

		// fetch the user by email
		user, err := h.models.Users.GetByEmail(input.Email)
		if err != nil {
			switch {
			case errors.Is(err, data.ErrRecordNotFound):
				responseWriter.NotFoundResponse(w, r)
			default:
				responseWriter.ServerErrorResponse(w, r, err)
			}
			return
		}

		// check if the password is correct
		match, err := user.Password.Matches(input.Password)
		if err != nil {
			responseWriter.ServerErrorResponse(w, r, err)
			return
		}
		if !match {
			responseWriter.UnauthorizedResponse(w, r, errors.New("invalid password"))
			return
		}

    // send the jwt token in the response
		jwt, err := jwt.GenerateJWT(jwtSecret, user.Email, user.ID)
		if err != nil {
			responseWriter.ServerErrorResponse(w, r, err)
			return
		}

		err = json.WriteJSON(json.Envelope{"token": jwt}, w, http.StatusOK, nil)
		if err != nil {
			responseWriter.ServerErrorResponse(w, r, err)
		}
	}
}
