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

type ModelHandler struct {
	Service yakit.ModelService
}

func (h ModelHandler) Model(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	model, err := h.Service.Model(vars["id"])
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error: %s", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(model)
}

func (h ModelHandler) Models(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()

	models, err := h.Service.Models(params.Get("brand"))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error: %s", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(models)
}

func (h ModelHandler) CreateModel(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Error: %s", err)
		return
	}

	var m yakit.Model
	err = json.Unmarshal(body, &m)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Error: %s", err)
		return
	}

	model, err := h.Service.CreateModel(m)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error: %s", err)
		return
	}

	w.Header().Set("Location", fmt.Sprintf("/models/%d", model.ID))
	w.WriteHeader(http.StatusCreated)
}

func (h ModelHandler) UpdateModel(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Error: %s", err)
		return
	}

	var m yakit.Model
	err = json.Unmarshal(body, &m)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Error: %s", err)
		return
	}

	vars := mux.Vars(r)

	m.ID, err = strconv.Atoi(vars["id"])

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Error: %s", err)
		return
	}

	model, err := h.Service.UpdateModel(m)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error: %s", err)
		return
	}

	w.Header().Set("Location", fmt.Sprintf("/models/%d", model.ID))
	w.WriteHeader(http.StatusOK)
}

func (h ModelHandler) DeleteModel(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	err := h.Service.DeleteModel(vars["id"])
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error: %s", err)
		return
	}

	w.WriteHeader(http.StatusOK)
}
