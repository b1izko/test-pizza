package order

import (
	"io/ioutil"
	"net/http"

	"github.com/b1izko/test-pizza/manager/internal/api"
	"github.com/b1izko/test-pizza/manager/store/order"
	"github.com/b1izko/test-pizza/manager/utils"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// GetOrder return order by ID
func (h *Handler) GetOrder(w http.ResponseWriter, r *http.Request) {

	_, err := utils.AuthToken(r, h.Store)
	if err != nil {
		api.WriteError(w, api.CodeAccessDenied, "Access Denied")
		return
	}

	value := mux.Vars(r)

	order, err := order.ByID(value["id"], h.Store)
	if err != nil {
		api.WriteError(w, api.CodeParseError, err.Error())
		return
	}
	response := api.NewResponse(order)
	response.WriteResponse(w)
	return
}

// CreateOrder create new order
func (h *Handler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil || len(body) == 0 {
		api.WriteError(w, api.CodeInvalidRequest, err.Error())
		return
	}

	defer r.Body.Close()

	var order order.Model
	err = order.UnmarshalJSON(body)
	if err != nil {
		api.WriteError(w, api.CodeParseError, err.Error())
		return
	}

	if err := order.Save(h.Store); err != nil {
		api.WriteError(w, api.CodeSaveError, err.Error())
		return
	}

	response := api.NewResponse(order)
	response.WriteResponse(w)
	return

}

// EditOrder edit resume by id
func (h *Handler) EditOrder(w http.ResponseWriter, r *http.Request) {
	_, err := utils.AuthToken(r, h.Store)
	if err != nil {
		api.WriteError(w, api.CodeAccessDenied, "Access Denied")
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil || len(body) == 0 {
		api.WriteError(w, api.CodeInvalidRequest, err.Error())
		return
	}

	defer r.Body.Close()

	var newOrder order.Model
	err = newOrder.UnmarshalJSON(body)
	if err != nil {
		api.WriteError(w, api.CodeParseError, err.Error())
		return
	}

	values := mux.Vars(r)

	id, err := primitive.ObjectIDFromHex(values["id"])
	if err != nil {
		api.WriteError(w, api.CodeParseError, err.Error())
		return
	}
	newOrder.ID = id

	if err := newOrder.Save(h.Store); err != nil {
		api.WriteError(w, api.CodeSaveError, err.Error())
		return
	}

	result, err := order.ByID(newOrder.ID.Hex(), h.Store)
	if err != nil {
		api.WriteError(w, api.CodeParseError, err.Error())
		return
	}

	response := api.NewResponse(result)
	response.WriteResponse(w)
	return
}

// RemoveOrder remove order
func (h *Handler) RemoveOrder(w http.ResponseWriter, r *http.Request) {

	_, err := utils.AuthToken(r, h.Store)
	if err != nil {
		api.WriteError(w, api.CodeAccessDenied, "Access Denied")
		return
	}

	order, err := order.ByID(mux.Vars(r)["id"], h.Store)
	if err != nil {
		api.WriteError(w, api.CodeParseError, err.Error())
		return
	}

	if err := order.Remove(h.Store); err != nil {
		api.WriteError(w, api.CodeDeleteError, err.Error())
		return
	}

	response := api.NewResponse(true)
	response.WriteResponse(w)
	return
}
