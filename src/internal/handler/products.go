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

type ProductHandler struct {
	q *db.Queries
}

func NewProductHandler(q *db.Queries) *ProductHandler {
	return &ProductHandler{q: q}
}

// List godoc
// @Summary     Список продуктов
// @Tags        products
// @Produce     json
// @Success     200  {array}   db.Product
// @Failure     500  {object}  map[string]string
// @Router      /products [get]
func (h *ProductHandler) List(w http.ResponseWriter, r *http.Request) {
	products, err := h.q.ListProducts(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if products == nil {
		products = []db.Product{}
	}
	writeJSON(w, http.StatusOK, products)
}

// Get godoc
// @Summary     Получить продукт
// @Tags        products
// @Produce     json
// @Param       id   path      string  true  "UUID продукта"
// @Success     200  {object}  db.Product
// @Failure     400  {object}  map[string]string
// @Failure     404  {object}  map[string]string
// @Failure     500  {object}  map[string]string
// @Router      /products/{id} [get]
func (h *ProductHandler) Get(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid uuid")
		return
	}
	product, err := h.q.GetProduct(r.Context(), id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			writeError(w, http.StatusNotFound, "product not found")
			return
		}
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, product)
}

// Create godoc
// @Summary     Создать продукт
// @Tags        products
// @Accept      json
// @Produce     json
// @Param       body  body      db.CreateProductParams  true  "Данные продукта"
// @Success     201   {object}  db.Product
// @Failure     400   {object}  map[string]string
// @Failure     500   {object}  map[string]string
// @Router      /products [post]
func (h *ProductHandler) Create(w http.ResponseWriter, r *http.Request) {
	var params db.CreateProductParams
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	params.ID = uuid.Must(uuid.NewV7())
	product, err := h.q.CreateProduct(r.Context(), params)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, product)
}

// Update godoc
// @Summary     Обновить продукт
// @Tags        products
// @Accept      json
// @Produce     json
// @Param       id    path      string                  true  "UUID продукта"
// @Param       body  body      db.UpdateProductParams  true  "Новые данные"
// @Success     200   {object}  db.Product
// @Failure     400   {object}  map[string]string
// @Failure     404   {object}  map[string]string
// @Failure     500   {object}  map[string]string
// @Router      /products/{id} [put]
func (h *ProductHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid uuid")
		return
	}
	var params db.UpdateProductParams
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	params.ID = id
	product, err := h.q.UpdateProduct(r.Context(), params)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			writeError(w, http.StatusNotFound, "product not found")
			return
		}
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, product)
}

// Delete godoc
// @Summary     Удалить продукт
// @Tags        products
// @Param       id   path  string  true  "UUID продукта"
// @Success     204
// @Failure     400  {object}  map[string]string
// @Failure     500  {object}  map[string]string
// @Router      /products/{id} [delete]
func (h *ProductHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid uuid")
		return
	}
	if err := h.q.DeleteProduct(r.Context(), id); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
