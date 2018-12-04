package server

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"github.com/egegunes/yakit/yakit"
)

type ModelHandler struct {
	service yakit.ModelService
	logger  *log.Logger
}

func NewModelHandler(s yakit.ModelService, l *log.Logger) *ModelHandler {
	return &ModelHandler{service: s, logger: l}
}

func (h ModelHandler) Model(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	model, err := h.service.Model(vars["id"])
	if err != nil {
		h.logger.Printf("couldn't get model: %v", err)
		http.Error(w, "couldn't get model", http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(model)
}

func (h ModelHandler) Models(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()

	models, err := h.service.Models(params.Get("brand_id"), params.Get("brand_name"))
	if err != nil {
		h.logger.Printf("couldn't get models: %v", err)
		http.Error(w, "couldn't get models", http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(models)
}

func (h ModelHandler) CreateModel(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		h.logger.Printf("couldn't read from request body: %v", err)
		http.Error(w, "couldn't read from request body", http.StatusBadRequest)
	}
	defer r.Body.Close()

	var m yakit.Model
	if err := json.Unmarshal(body, &m); err != nil {
		h.logger.Printf("couldn't unmarshal json: %v", err)
		http.Error(w, "couldn't unmarshal json", http.StatusBadRequest)
	}

	model, err := h.service.CreateModel(m)
	if err != nil {
		h.logger.Printf("couldn't create model: %v", err)
		http.Error(w, "couldn't create model", http.StatusBadRequest)
	}

	w.Header().Set("Location", fmt.Sprintf("/models/%d", model.ID))
	w.WriteHeader(http.StatusCreated)
}

func (h ModelHandler) UpdateModel(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		h.logger.Printf("couldn't read from request body: %v", err)
		http.Error(w, "couldn't read from request body", http.StatusBadRequest)
	}
	defer r.Body.Close()

	var m yakit.Model
	if err := json.Unmarshal(body, &m); err != nil {
		h.logger.Printf("couldn't unmarshal json: %v", err)
		http.Error(w, "couldn't unmarshal json", http.StatusBadRequest)
	}

	vars := mux.Vars(r)

	m.ID, err = strconv.Atoi(vars["id"])
	if err != nil {
		h.logger.Printf("couldn't convert ID %s to integer: %v", m.ID, err)
		http.Error(w, "couldn't convert ID to integer", http.StatusBadRequest)
	}

	model, err := h.service.UpdateModel(m)
	if err != nil {
		h.logger.Printf("couldn't update model: %v", err)
		http.Error(w, "couldn't update model", http.StatusInternalServerError)
	}

	w.Header().Set("Location", fmt.Sprintf("/models/%d", model.ID))
	w.WriteHeader(http.StatusOK)
}

func (h ModelHandler) DeleteModel(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	if err := h.service.DeleteModel(vars["id"]); err != nil {
		h.logger.Printf("couldn't delete model: %v", err)
		http.Error(w, "couldn't delete model", http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusOK)
}
