package user

import (
	"net/http"
)

type Handler struct{}

func NewHandler() *Handler{
  return &Handler{}
}

func (h *Handler) RegisterUserHandler(w http.ResponseWriter, r *http.Request){
  // 
}

func (h *Handler) LoginUserHandler(w http.ResponseWriter, r *http.Request){
  // 
}
