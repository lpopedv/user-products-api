package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/lpopedv/user-products-api/infra/database"
	"github.com/lpopedv/user-products-api/internal/dto"
	"github.com/lpopedv/user-products-api/internal/entity"
)

type ProductHandler struct {
	ProductDB database.ProductInterface
}

func NewProductHandler(db database.ProductInterface) *ProductHandler {
	return &ProductHandler{
		ProductDB: db,
	}
}

func (h *ProductHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	var product dto.CreateProductInput

	err := json.NewDecoder(r.Body).Decode(&product)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	p, err := entity.NewProduct(product.Name, product.Price)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = h.ProductDB.Create(p)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}
