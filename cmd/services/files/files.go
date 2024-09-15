package files

import (
	encjson "encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"sarath/backend_project/internal/cache"
	"sarath/backend_project/internal/data"
	filestore "sarath/backend_project/internal/file_store"
	"sarath/backend_project/internal/json"
	"sarath/backend_project/internal/response"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type Handler struct {
	Logger    *log.Logger
	models    *data.Models
	fileStore *filestore.FileStore
	Cache     *cache.Cache
}

func NewHandler(logger *log.Logger, models *data.Models, fileStore *filestore.FileStore, cache *cache.Cache) *Handler {
	return &Handler{
		Logger:    logger,
		models:    models,
		fileStore: fileStore,
		Cache:     cache,
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
	fileURL, err := h.fileStore.UploadFile(file, name)
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

func (h *Handler) ShareFileHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	responseWriter := response.NewResponseWriter(h.Logger)
	fileID, err := strconv.ParseInt(vars["file_id"], 10, 64)
	if err != nil {
		responseWriter.BadRequestResponse(w, r, err)
		return
	}

	// check if the file url is in cache
	fileUrl, err := h.Cache.Get(strconv.FormatInt(fileID, 10))
	if err == nil {
		data := json.Envelope{"file_url": fileUrl}
		err = json.WriteJSON(data, w, http.StatusOK, nil)
		if err != nil {
			responseWriter.ServerErrorResponse(w, r, err)
		}
		return
	}

	// return the file url whose id is fileID
	metadata, err := h.models.MetaData.Get(fileID)
	if err != nil {
		responseWriter.ServerErrorResponse(w, r, err)
		return
	}

	// write metadata to cache
	err = h.Cache.Set(strconv.FormatInt(fileID, 10), metadata.FileUrl, time.Minute*5)
	if err != nil {
		h.Logger.Printf("error writing to cache: %v", err)
	}

	data := json.Envelope{"file_url": metadata.FileUrl}
	err = json.WriteJSON(data, w, http.StatusOK, nil)
	if err != nil {
		responseWriter.ServerErrorResponse(w, r, err)
	}
}

func (h *Handler) GetFilesHandler(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value("id").(int64)
	cacheKey := strconv.FormatInt(id, 10)
	responseWriter := response.NewResponseWriter(h.Logger)

	// check if the metadata is in cache
	cacheMetaData, err := h.Cache.Get(cacheKey)
	if err == nil {
		var metadataList []*data.MetaData
		err = encjson.Unmarshal([]byte(cacheMetaData), &metadataList)
		if err != nil {
			h.Logger.Printf("error unmarshalling metadata: %v", err)
		}
		data := json.Envelope{"metadata": metadataList}
		err = json.WriteJSON(data, w, http.StatusOK, nil)
		if err != nil {
			responseWriter.ServerErrorResponse(w, r, err)
		}
		return
	}

	// TODO: Handle pagination
	// get the files based on the user id
	metadata, err := h.models.MetaData.GetByUserID(id)
	if err != nil {
		responseWriter.ServerErrorResponse(w, r, err)
		return
	}

	// write metadata to cache
	jsonData, err := encjson.Marshal(metadata)
	if err != nil {
		h.Logger.Printf("error marshalling metadata: %v", err)
	}
	err = h.Cache.Set(cacheKey, string(jsonData), time.Minute*5)
	if err != nil {
		h.Logger.Printf("error writing to cache: %v", err)
	}

	// send the metadata to the client
	data := json.Envelope{"metadata": metadata}
	err = json.WriteJSON(data, w, http.StatusOK, nil)
	if err != nil {
		responseWriter.ServerErrorResponse(w, r, err)
	}
}

func (h *Handler) SearchFileHandler(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value("id").(int64)

	// Retrieve query parameters
	params := r.URL.Query()
	filename := params.Get("filename")
	contentType := params.Get("content_type")
	timeString := params.Get("time")

	responseWriter := response.NewResponseWriter(h.Logger)

	// convert the time string to time.Time
	var parsedTime time.Time
	var err error
	if timeString != "" {
		parsedTime, err = time.Parse(time.RFC3339, timeString)
		if err != nil {
			responseWriter.BadRequestResponse(w, r, fmt.Errorf("invalid time format time should be in RFC3339 format"))
			return
		}
	}

	// search the files based on the input
	metadata, err := h.models.MetaData.Search(id, filename, contentType, &parsedTime)
	if err != nil {
		responseWriter.ServerErrorResponse(w, r, err)
		return
	}

	data := json.Envelope{"metadata": metadata}
	err = json.WriteJSON(data, w, http.StatusOK, nil)
	if err != nil {
		responseWriter.ServerErrorResponse(w, r, err)
	}
}
