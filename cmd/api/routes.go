package api

import (
	"net/http"
	"sarath/backend_project/cmd/services/files"
	"sarath/backend_project/cmd/services/user"

	"github.com/gorilla/mux"
)

func (app *Application) Routes() *mux.Router {
	router := mux.NewRouter()

	subrouter := router.PathPrefix("/api/v1").Subrouter()
  subrouter.Use(getAuthMiddlewarewithJWT(app.Config.Jwt.Secret))

  userHandler := user.NewHandler(app.Logger, app.Models);
  filesHandler := files.NewHandler(app.Logger, app.Models, app.S3Sess, app.Config.Aws.Bucket);

  router.HandleFunc("/register", userHandler.RegisterUserHandler).Methods(http.MethodPost)
  router.HandleFunc("/login", userHandler.GetLoginUserHandler(app.Config.Jwt.Secret)).Methods(http.MethodPost)

  // TODO: Keep pagination, caching in mind
  subrouter.HandleFunc("/upload", filesHandler.UploadFileHandler).Methods(http.MethodPost)
  subrouter.HandleFunc("/share/{file_id}", filesHandler.ShareFileHandler).Methods(http.MethodGet)
  subrouter.HandleFunc("/files", filesHandler.GetFilesHandler).Methods(http.MethodGet)
  subrouter.HandleFunc("/search", filesHandler.SearchFileHandler).Methods(http.MethodGet)

	return router
}
