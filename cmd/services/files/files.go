package files

import (
	"fmt"
	"log"
	"mime/multipart"
	"net/http"

	"sarath/backend_project/internal/data"
	"sarath/backend_project/internal/json"
	"sarath/backend_project/internal/response"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/google/uuid"
)

type Handler struct {
	Logger *log.Logger
	S3Sess *session.Session
	models *data.Models
	Bucket string
}

func NewHandler(logger *log.Logger, models *data.Models, s3Sess *session.Session, bucket string) *Handler {
	return &Handler{
		Logger: logger,
		S3Sess: s3Sess,
		models: models,
		Bucket: bucket,
	}
}

func (h *Handler) UploadFileHandler(w http.ResponseWriter, r *http.Request) {
	responseWriter := response.NewResponseWriter(h.Logger)

	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		responseWriter.FileTooLargeResponse(w, r)
		return
	}

	file, handler, err := r.FormFile("file")
	if err != nil {
		responseWriter.BadRequestResponse(w, r, err)
		return
	}
	defer file.Close()

	// upload the file to s3 with a unique name
	name := fmt.Sprintf("%s-%s", uuid.New().String(), handler.Filename)

	// upload the file to s3
	fileURL, err := h.uploadToS3(file, name)
	if err != nil {
		responseWriter.ServerErrorResponse(w, r, err)
		return
	}

	// get user id from context
	id := r.Context().Value("id").(int64)

	// create the metadata object
	metadata := &data.MetaData{
		UserId:      id,
		Name:        handler.Filename,
		Size:        handler.Size,
		ContentType: handler.Header.Get("Content-Type"),
		FileUrl:     fileURL,
	}

	// insert the metadata into the database
	err = h.models.MetaData.Insert(metadata)
	if err != nil {
		responseWriter.ServerErrorResponse(w, r, err)
		return
	}

	// send the file url and metadata to the client
	data := json.Envelope{"file_url": fileURL, "metadata": metadata}
	err = json.WriteJSON(data, w, http.StatusCreated, nil)
	if err != nil {
		responseWriter.ServerErrorResponse(w, r, err)
	}
}

func (h *Handler) uploadToS3(file multipart.File, s3Key string) (string, error) {
	uploader := s3manager.NewUploader(h.S3Sess)

	result, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(h.Bucket),
		Key:    aws.String(s3Key),
		Body:   file,
	})
	if err != nil {
		return "", err
	}

	return result.Location, nil
}
