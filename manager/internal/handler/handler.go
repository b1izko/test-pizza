package handler

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime/debug"

	"github.com/b1izko/test-pizza/manager/internal/api/order"
	"github.com/b1izko/test-pizza/manager/store"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

var ctx = context.TODO()

// Handler for this app
type Handler struct {
	store  store.Storage
	router *mux.Router
}

// New created handler for given storage and init all routes
func New(store store.Storage) *Handler {
	h := &Handler{
		store:  store,
		router: mux.NewRouter(),
	}
	orderHandler := order.Handler{Store: store}
	orderHandler.Get(h.router)

	os.MkdirAll("uploads", 0655)

	return h
}

// Get http.Handler
func (h *Handler) Get() http.Handler {
	return panicMiddleware(handlers.CORS(
		handlers.AllowedHeaders([]string{"Content-Type", "Authorization"}),
		handlers.AllowedOrigins([]string{"*"}),
		handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"}),
	)(
		h.requestOptionsMiddleware(h.requestLogMiddleware(h.router)),
	))
}

// Close handler storage
func (h *Handler) Close() error {
	return h.store.Disconnect()
}

func panicMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				fmt.Println(`\e[33m[ PANIC ]\e[0m :: `, err, string(debug.Stack()))
				http.Error(w,
					http.StatusText(http.StatusInternalServerError),
					http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func (h *Handler) requestLogMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Print(map[string]interface{}{
			"Method": r.Method,
			"URL":    r.URL,
			"Agent":  r.Header.Get("User-Agent"),
		})
		next.ServeHTTP(w, r)
	})
}

func (h *Handler) requestOptionsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "OPTIONS" {
			return
		}
		next.ServeHTTP(w, r)
	})
}
