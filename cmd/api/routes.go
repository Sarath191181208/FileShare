package api

import (
	"net/http"

	"github.com/gorilla/mux"
)

func (app *Application) Routes() *mux.Router {
	router := mux.NewRouter()
	subrouter := router.PathPrefix("/api/v1").Subrouter()
  subrouter.Use(getAuthMiddlewarewithJWT(app.Config.Jwt.Secret))

	// router.HandleFunc("/register", app.registerUserHandler).Methods(http.MethodPost)
	// router.HandleFunc("/login", app.loginUserHandler).Methods(http.MethodPost)

  // TODO: Keep pagination, caching in mind
	// subrouter.HandleFunc("/files", app.getFilesMetadata).Methods(http.MethodGet)
  // subrouter.HandleFunc("/share/:file_id", app.createMovieHandler).Methods(http.MethodPost)
  // subrouter.HandleFunc("/search", app.searchMoviesHandler).Methods(http.MethodGet)
	// subrouter.HandleFunc("/v1/movies/:id", app.showMovieHandler).Methods(http.MethodGet)

	return router
}
