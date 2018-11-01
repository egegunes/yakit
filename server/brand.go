package server

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"yakit/yakit"
)

type BrandHandler struct {
	Service yakit.BrandService
}

func (h BrandHandler) Brand(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	brand, err := h.Service.Brand(vars["id"])
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error: %s", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(brand)
}

func (h BrandHandler) Brands(w http.ResponseWriter, r *http.Request) {
	brands, err := h.Service.Brands()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error: %s", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(brands)
}

func (h BrandHandler) CreateBrand(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Error: %s", err)
		return
	}

	var b yakit.Brand
	err = json.Unmarshal(body, &b)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Error: %s", err)
		return
	}

	brand, err := h.Service.CreateBrand(b)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error: %s", err)
		return
	}

	w.Header().Set("Location", fmt.Sprintf("/brand/%d", brand.ID))
	w.WriteHeader(http.StatusCreated)
}

func (h BrandHandler) UpdateBrand(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Error: %s", err)
		return
	}

	var b yakit.Brand
	err = json.Unmarshal(body, &b)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Error: %s", err)
		return
	}

	vars := mux.Vars(r)

	b.ID, err = strconv.Atoi(vars["id"])

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Error: %s", err)
		return
	}

	brand, err := h.Service.UpdateBrand(b)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error: %s", err)
		return
	}

	w.Header().Set("Location", fmt.Sprintf("/brand/%d", brand.ID))
	w.WriteHeader(http.StatusOK)
}

func (h BrandHandler) DeleteBrand(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	err := h.Service.DeleteBrand(vars["id"])
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error: %s", err)
		return
	}

	w.WriteHeader(http.StatusOK)
}
