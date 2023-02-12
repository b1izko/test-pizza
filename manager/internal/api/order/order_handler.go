package order

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

	router.HandleFunc("/order/create", h.CreateOrder).Methods("POST")

	router.HandleFunc("/order/edit", h.EditOrder).
		Queries("id", "{id}").
		Methods("POST")

	router.HandleFunc("/order/remove", h.RemoveOrder).
		Queries("id", "{id}").
		Methods("GET")

	router.HandleFunc("/order/get", h.GetOrder).
		Queries("id", "{id}").
		Methods("GET")

}
