package files

import (
	"log"
	"net/http"
)

type Handler struct {
	Logger *log.Logger
}

func NewHandler(logger *log.Logger) *Handler {
	return &Handler{logger}
}


func (h *Handler) UploadFileHandler(w http.ResponseWriter, r *http.Request) {
}
