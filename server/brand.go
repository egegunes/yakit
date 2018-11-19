package server

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"yakit/yakit"
)

type BrandHandler struct {
	logger  *log.Logger
	service yakit.BrandService
}

func NewBrandHandler(s yakit.BrandService, l *log.Logger) *BrandHandler {
	return &BrandHandler{service: s, logger: l}
}

func (h BrandHandler) Brand(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	brand, err := h.service.Brand(vars["id"])
	if err != nil {
		h.logger.Printf("could not get brand: %v", err)
		http.Error(w, "could not get brand", http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(brand); err != nil {
		h.logger.Printf("could not encode brand: %v", err)
		http.Error(w, "could not encode brand", http.StatusInternalServerError)
	}
}

func (h BrandHandler) Brands(w http.ResponseWriter, r *http.Request) {
	brands, err := h.service.Brands()
	if err != nil {
		h.logger.Printf("could not get brands: %v", err)
		http.Error(w, "could not get brands", http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(brands)
}

func (h BrandHandler) CreateBrand(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		h.logger.Printf("could not read from body: %v", err)
		http.Error(w, "could not read from body", http.StatusBadRequest)
	}

	var b yakit.Brand
	if err := json.Unmarshal(body, &b); err != nil {
		h.logger.Printf("could not unmarshal json: %v", err)
		http.Error(w, "could not unmarshal json", http.StatusBadRequest)
	}

	brand, err := h.service.CreateBrand(b)
	if err != nil {
		h.logger.Printf("could not create brand: %v", err)
		http.Error(w, "could not create brand", http.StatusBadRequest)
	}

	w.Header().Set("Location", fmt.Sprintf("/brand/%d", brand.ID))
	w.WriteHeader(http.StatusCreated)
}

func (h BrandHandler) UpdateBrand(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		h.logger.Printf("could not read from body: %v", err)
		http.Error(w, "could not read from body", http.StatusBadRequest)
	}

	var b yakit.Brand
	if err := json.Unmarshal(body, &b); err != nil {
		h.logger.Printf("could not unmarshal json: %v", err)
		http.Error(w, "could not unmarshal json", http.StatusBadRequest)
	}

	vars := mux.Vars(r)

	b.ID, err = strconv.Atoi(vars["id"])
	if err != nil {
		h.logger.Printf("could not convert id to string: %v", err)
		http.Error(w, "could not convert id to string: %v", http.StatusBadRequest)
	}

	brand, err := h.service.UpdateBrand(b)
	if err != nil {
		h.logger.Printf("could not update brand: %v", err)
		http.Error(w, "could not update brand", http.StatusInternalServerError)
	}

	w.Header().Set("Location", fmt.Sprintf("/brand/%d", brand.ID))
	w.WriteHeader(http.StatusOK)
}

func (h BrandHandler) DeleteBrand(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	err := h.service.DeleteBrand(vars["id"])
	if err != nil {
		h.logger.Printf("could not delete brand: %v", err)
		http.Error(w, "could not delete brand", http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusOK)
}
