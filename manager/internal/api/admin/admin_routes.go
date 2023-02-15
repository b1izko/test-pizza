package admin

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/b1izko/test-pizza/manager/internal/api"
	"github.com/b1izko/test-pizza/manager/store/admin"
	"github.com/b1izko/test-pizza/manager/utils"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type authData struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

// Auth return JWT for request
func (h *Handler) Auth(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		err := api.Error{
			Code: api.CodeGeneralError,
			Msg:  "Internal server error",
		}
		response := err.GetResponse()
		response.WriteResponse(w)
		return
	}

	auth := &authData{}
	err = json.Unmarshal(body, auth)
	if err != nil {
		err := api.Error{
			Code: api.CodeGeneralError,
			Msg:  "Internal server error",
		}
		response := err.GetResponse()
		response.WriteResponse(w)
		return
	}

	data, err := admin.Check(auth.Login, utils.Hash(auth.Password), h.Store)
	if err != nil {
		_err := api.Error{
			Code: api.CodeAccessDenied,
			Msg:  err.Error(),
		}
		response := _err.GetResponse()
		response.WriteResponse(w)
		return
	}

	user, err := admin.ByID(data.ID.Hex(), h.Store)
	if err != nil {
		err := api.Error{
			Code: api.CodeInvalidUser,
			Msg:  "Internal server error",
		}
		response := err.GetResponse()
		response.WriteResponse(w)
		return
	}

	token, err := user.GetJWT()
	if err != nil {
		err := api.Error{
			Code: api.CodeInvalidToken,
			Msg:  "Internal server error",
		}
		response := err.GetResponse()
		response.WriteResponse(w)
		return
	}

	user.LastAuth = time.Now()
	user.Save(h.Store)

	response := api.NewResponse(token)
	response.WriteResponse(w)
}

// GetAdmin return admin by ID
func (h *Handler) GetAdmin(w http.ResponseWriter, r *http.Request) {

	_, err := utils.AuthToken(r, h.Store)
	if err != nil {
		api.WriteError(w, api.CodeAccessDenied, "Access Denied")
		return
	}

	value := mux.Vars(r)

	admin, err := admin.ByID(value["id"], h.Store)
	if err != nil {
		api.WriteError(w, api.CodeParseError, err.Error())
		return
	}
	response := api.NewResponse(admin)
	response.WriteResponse(w)
	return
}

// CreateAdmin create new admin
func (h *Handler) CreateAdmin(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil || len(body) == 0 {
		api.WriteError(w, api.CodeInvalidRequest, err.Error())
		return
	}

	defer r.Body.Close()

	var admin admin.Model
	err = admin.UnmarshalJSON(body)
	if err != nil {
		api.WriteError(w, api.CodeParseError, err.Error())
		return
	}

	if err := admin.Save(h.Store); err != nil {
		api.WriteError(w, api.CodeSaveError, err.Error())
		return
	}

	response := api.NewResponse(admin)
	response.WriteResponse(w)
	return

}

// EditAdmin edit admin by id
func (h *Handler) EditAdmin(w http.ResponseWriter, r *http.Request) {
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

	var newAdmin admin.Model
	err = newAdmin.UnmarshalJSON(body)
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
	newAdmin.ID = id

	if err := newAdmin.Save(h.Store); err != nil {
		api.WriteError(w, api.CodeSaveError, err.Error())
		return
	}

	result, err := admin.ByID(newAdmin.ID.Hex(), h.Store)
	if err != nil {
		api.WriteError(w, api.CodeParseError, err.Error())
		return
	}

	response := api.NewResponse(result)
	response.WriteResponse(w)
	return
}

// RemoveAdmin remove admin
func (h *Handler) RemoveAdmin(w http.ResponseWriter, r *http.Request) {

	_, err := utils.AuthToken(r, h.Store)
	if err != nil {
		api.WriteError(w, api.CodeAccessDenied, "Access Denied")
		return
	}

	admin, err := admin.ByID(mux.Vars(r)["id"], h.Store)
	if err != nil {
		api.WriteError(w, api.CodeParseError, err.Error())
		return
	}

	if err := admin.Remove(h.Store); err != nil {
		api.WriteError(w, api.CodeDeleteError, err.Error())
		return
	}

	response := api.NewResponse(true)
	response.WriteResponse(w)
	return
}
