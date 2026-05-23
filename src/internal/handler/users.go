package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"go-crud/internal/db"
)

type UserHandler struct {
	q *db.Queries
}

func NewUserHandler(q *db.Queries) *UserHandler {
	return &UserHandler{q: q}
}

// List godoc
// @Summary     Список пользователей
// @Tags        users
// @Produce     json
// @Success     200  {array}   db.User
// @Failure     500  {object}  map[string]string
// @Router      /users [get]
func (h *UserHandler) List(w http.ResponseWriter, r *http.Request) {
	users, err := h.q.ListUsers(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if users == nil {
		users = []db.User{}
	}
	writeJSON(w, http.StatusOK, users)
}

// Get godoc
// @Summary     Получить пользователя
// @Tags        users
// @Produce     json
// @Param       id   path      string  true  "UUID пользователя"
// @Success     200  {object}  db.User
// @Failure     400  {object}  map[string]string
// @Failure     404  {object}  map[string]string
// @Failure     500  {object}  map[string]string
// @Router      /users/{id} [get]
func (h *UserHandler) Get(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid uuid")
		return
	}
	user, err := h.q.GetUser(r.Context(), id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			writeError(w, http.StatusNotFound, "user not found")
			return
		}
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, user)
}

// Create godoc
// @Summary     Создать пользователя
// @Tags        users
// @Accept      json
// @Produce     json
// @Param       body  body      db.CreateUserParams  true  "Данные пользователя"
// @Success     201   {object}  db.User
// @Failure     400   {object}  map[string]string
// @Failure     500   {object}  map[string]string
// @Router      /users [post]
func (h *UserHandler) Create(w http.ResponseWriter, r *http.Request) {
	var params db.CreateUserParams
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	params.ID = uuid.Must(uuid.NewV7())
	user, err := h.q.CreateUser(r.Context(), params)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, user)
}

// Update godoc
// @Summary     Обновить пользователя
// @Tags        users
// @Accept      json
// @Produce     json
// @Param       id    path      string               true  "UUID пользователя"
// @Param       body  body      db.UpdateUserParams  true  "Новые данные"
// @Success     200   {object}  db.User
// @Failure     400   {object}  map[string]string
// @Failure     404   {object}  map[string]string
// @Failure     500   {object}  map[string]string
// @Router      /users/{id} [put]
func (h *UserHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid uuid")
		return
	}
	var params db.UpdateUserParams
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	params.ID = id
	user, err := h.q.UpdateUser(r.Context(), params)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			writeError(w, http.StatusNotFound, "user not found")
			return
		}
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, user)
}

// Delete godoc
// @Summary     Удалить пользователя
// @Tags        users
// @Param       id   path  string  true  "UUID пользователя"
// @Success     204
// @Failure     400  {object}  map[string]string
// @Failure     500  {object}  map[string]string
// @Router      /users/{id} [delete]
func (h *UserHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid uuid")
		return
	}
	if err := h.q.DeleteUser(r.Context(), id); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
