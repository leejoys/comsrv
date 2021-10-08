package api

import (
	"comsrv/pkg/storage"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

// API of comments server
type API struct {
	db storage.Interface
	r  *mux.Router
}

// API constructor
func New(db storage.Interface) *API {
	api := API{
		db: db,
	}
	api.r = mux.NewRouter()
	api.endpoints()
	return &api
}

// API handler registration
func (api *API) endpoints() {
	// получить комментарии к новости n
	api.r.HandleFunc("/comments", api.comments).Methods(http.MethodGet)
	// сохранить комментарий
	api.r.HandleFunc("/comments/{parentpost}/{parentcomment}", api.storeComment).Methods(http.MethodPost)

}

// Получение маршрутизатора запросов.
// Требуется для передачи маршрутизатора веб-серверу.
func (api *API) Router() *mux.Router {
	return api.r
}

// Получение всех публикаций.
func (api *API) comments(w http.ResponseWriter, r *http.Request) {
	ns := mux.Vars(r)["n"]
	n, err := strconv.Atoi(ns)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	posts, err := api.db.CommentsN(n)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	bytes, err := json.Marshal(posts)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(bytes)
}

// сохранение комментария
func (api *API) storeComment(w http.ResponseWriter, r *http.Request) {
	var err error
	c := storage.Comment{}
	c.ParentPost, err = strconv.Atoi(mux.Vars(r)["parentpost"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	c.ParentComment, err = strconv.Atoi(mux.Vars(r)["parentcomment"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	c.PubTime = time.Now().Unix()

	err = api.db.AddComment(c)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
