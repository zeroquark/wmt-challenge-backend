package handlers

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
	"wmt-challenge/db"
	"wmt-challenge/util"
)

type ProductsHandler struct {
	l *log.Logger
}

func NewProductsHandler(l *log.Logger) *ProductsHandler {
	return &ProductsHandler{l}
}

func (ph *ProductsHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	dbx := db.GetDB()
	dbx.Close()
}

// Products by Id
func (ph *ProductsHandler) GetProductById(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		ph.l.Println(err)
		util.RespondWithError(rw, http.StatusBadRequest, "Invalid product ID")
		return
	}

	dbx := db.GetDB()
	product, err := dbx.FindById(id)
	if err != nil {
		ph.l.Println(err)
		util.RespondWithError(rw, http.StatusInternalServerError, err.Error())
		return
	}
	if product == nil {
		ph.l.Println(err)
		util.RespondWithError(rw, http.StatusNotFound, "Product not found")
		return
	}

	if util.IsPalindrome(strconv.Itoa(id)) {
		product.Price = product.Price / 2
	}

	util.RespondWithJSON(rw, http.StatusOK, product)
}

// Products by Brand and/or Description
func (ph *ProductsHandler) GetProductByToken(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	token, present := vars["token"]
	if !present {
		ph.l.Println("Empty input")
		util.RespondWithError(rw, http.StatusBadRequest, "Empty input")
		return
	}

	var productsByBrand db.Products
	var productsByDescription db.Products
	var errBrand, errDescription error

	dbx := db.GetDB()
	productsByBrand, errBrand = dbx.FindByBrand(token)
	productsByDescription, errDescription = dbx.FindByDescription(token)
	if errBrand != nil || errDescription != nil {
		ph.l.Println(errBrand)
		util.RespondWithError(rw, http.StatusInternalServerError, errBrand.Error())
		return
	}
	if len(productsByBrand) == 0 && len(productsByDescription) == 0 {
		ph.l.Println(errBrand)
		util.RespondWithError(rw, http.StatusNotFound, "No product found for such token")
		return
	}

	products := append(productsByBrand, productsByDescription...)
	// products := productsByDescription
	_ = products.ToJSON(rw)
}
