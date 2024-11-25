package api

import (
	"encoding/json"
	"go-http-crud/database"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/google/uuid"
)

type Response struct {
	Data any `json:"data,omitempty"`
	Error string `json:"error,omitempty"`
}

func sendJSON(w http.ResponseWriter, r Response, status int) {
	w.Header().Set("Content-Type", "application/json")

	data, err := json.Marshal(r)
	if err != nil {
		slog.Error("failed to marshal json data", "error", err)
		sendJSON(w, Response{Error: "something went wrong"}, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(status)
	if _, err := w.Write(data); err != nil {
		slog.Error("failed to write response to client", "error", err)
		return
	}
}

func NewHandler(application database.Application) http.Handler {
	r := chi.NewMux()

	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)

	r.Post("/api/users", handleCreateUser(application))
	r.Get("/api/users", handleFetchAllUsers(application))
	r.Get("/api/users/{id}", handleFindUserById(application))
	r.Delete("/api/users/{id}", handleDeleteUser(application))
	r.Put("/api/users/{id}", handleUpdateUser(application))

	return r
}

func handleCreateUser(application database.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var body database.User
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			sendJSON(w, Response{Error: "unprocessable entity"}, http.StatusUnprocessableEntity)
			return
		}
		application.Insert(body)
		sendJSON(w, Response{Data: "user created"}, http.StatusOK)
	}
}

func handleFetchAllUsers(application database.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		users := application.FindAll()
		sendJSON(w, Response{Data: users}, http.StatusOK)
	}
}

func handleFindUserById(application database.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		user := application.FindById(database.Id(uuid.MustParse(id)))
		if user == nil {
			sendJSON(w, Response{Error: "user not found"}, http.StatusNotFound)
			return
		}
		sendJSON(w, Response{Data: user}, http.StatusOK)
	}
}

func handleDeleteUser(application database.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		parsedID := database.Id(uuid.MustParse(id))
		application.Delete(parsedID)
	}
}

func handleUpdateUser(application database.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		parsedID := database.Id(uuid.MustParse(id))

		var body database.User
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			sendJSON(w, Response{Error: "unprocessable body"}, http.StatusUnprocessableEntity)
			return
		}

		application.Update(parsedID, body)
	}
}