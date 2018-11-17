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

	models, err := h.service.Models(params.Get("brand"))
	if err != nil {
		h.logger.Printf("couldn't get models: %v", err)
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
		h.logger.Printf("couldn't read from request body: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Error: %s", err)
		return
	}
	defer r.Body.Close()

	var m yakit.Model
	if err := json.Unmarshal(body, &m); err != nil {
		h.logger.Printf("couldn't unmarshal json to Model: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Error: %s", err)
		return
	}

	model, err := h.service.CreateModel(m)
	if err != nil {
		h.logger.Printf("couldn't create Model: %v", err)
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
		h.logger.Printf("couldn't read from request body: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Error: %s", err)
		return
	}
	defer r.Body.Close()

	var m yakit.Model
	if err := json.Unmarshal(body, &m); err != nil {
		h.logger.Printf("couldn't unmarshal json to Model: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Error: %s", err)
		return
	}

	vars := mux.Vars(r)

	m.ID, err = strconv.Atoi(vars["id"])
	if err != nil {
		h.logger.Printf("couldn't convert ID %s to integer: %v", m.ID, err)
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Error: %s", err)
		return
	}

	model, err := h.service.UpdateModel(m)
	if err != nil {
		h.logger.Printf("couldn't update Model: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error: %s", err)
		return
	}

	w.Header().Set("Location", fmt.Sprintf("/models/%d", model.ID))
	w.WriteHeader(http.StatusOK)
}

func (h ModelHandler) DeleteModel(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	if err := h.service.DeleteModel(vars["id"]); err != nil {
		h.logger.Printf("couldn't delete Model: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error: %s", err)
		return
	}

	w.WriteHeader(http.StatusOK)
}
