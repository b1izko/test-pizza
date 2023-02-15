package admin

import (
	"github.com/b1izko/test-pizza/manager/store"

	"github.com/gorilla/mux"
)

// Handler ...
type Handler struct {
	Store store.Storage
}

// Get is get
func (h *Handler) Get(router *mux.Router) {

	router.HandleFunc("/admin/create", h.CreateAdmin).Methods("POST")

	router.HandleFunc("/admin/edit", h.EditAdmin).
		Queries("id", "{id}").
		Methods("POST")

	router.HandleFunc("/admin/remove", h.RemoveAdmin).
		Queries("id", "{id}").
		Methods("GET")

	router.HandleFunc("/admin/get", h.GetAdmin).
		Queries("id", "{id}").
		Methods("GET")

}
