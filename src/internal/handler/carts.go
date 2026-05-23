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

type CartHandler struct {
	q *db.Queries
}

func NewCartHandler(q *db.Queries) *CartHandler {
	return &CartHandler{q: q}
}

// ListByUser godoc
// @Summary     Корзина пользователя
// @Tags        carts
// @Produce     json
// @Param       userID  path      string  true  "UUID пользователя"
// @Success     200     {array}   db.Cart
// @Failure     400     {object}  map[string]string
// @Failure     500     {object}  map[string]string
// @Router      /carts/user/{userID} [get]
func (h *CartHandler) ListByUser(w http.ResponseWriter, r *http.Request) {
	userID, err := uuid.Parse(chi.URLParam(r, "userID"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid uuid")
		return
	}
	items, err := h.q.ListCartByUser(r.Context(), userID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if items == nil {
		items = []db.Cart{}
	}
	writeJSON(w, http.StatusOK, items)
}

// Add godoc
// @Summary     Добавить товар в корзину
// @Tags        carts
// @Accept      json
// @Produce     json
// @Param       body  body      db.AddToCartParams  true  "Товар и количество"
// @Success     201   {object}  db.Cart
// @Failure     400   {object}  map[string]string
// @Failure     500   {object}  map[string]string
// @Router      /carts [post]
func (h *CartHandler) Add(w http.ResponseWriter, r *http.Request) {
	var params db.AddToCartParams
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	params.ID = uuid.Must(uuid.NewV7())
	item, err := h.q.AddToCart(r.Context(), params)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, item)
}

// Update godoc
// @Summary     Изменить количество товара в корзине
// @Tags        carts
// @Accept      json
// @Produce     json
// @Param       id    path      string                   true  "UUID записи корзины"
// @Param       body  body      db.UpdateCartItemParams  true  "Новое количество"
// @Success     200   {object}  db.Cart
// @Failure     400   {object}  map[string]string
// @Failure     404   {object}  map[string]string
// @Failure     500   {object}  map[string]string
// @Router      /carts/{id} [put]
func (h *CartHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid uuid")
		return
	}
	var params db.UpdateCartItemParams
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	params.ID = id
	item, err := h.q.UpdateCartItem(r.Context(), params)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			writeError(w, http.StatusNotFound, "cart item not found")
			return
		}
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, item)
}

// Remove godoc
// @Summary     Удалить товар из корзины
// @Tags        carts
// @Param       id   path  string  true  "UUID записи корзины"
// @Success     204
// @Failure     400  {object}  map[string]string
// @Failure     500  {object}  map[string]string
// @Router      /carts/{id} [delete]
func (h *CartHandler) Remove(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid uuid")
		return
	}
	if err := h.q.RemoveFromCart(r.Context(), id); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
