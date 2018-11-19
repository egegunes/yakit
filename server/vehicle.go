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

type VehicleHandler struct {
	logger  *log.Logger
	service yakit.VehicleService
}

func NewVehicleHandler(s yakit.VehicleService, l *log.Logger) *VehicleHandler {
	return &VehicleHandler{service: s, logger: l}
}

func (h VehicleHandler) Vehicle(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	vehicle, err := h.service.Vehicle(vars["id"])
	if err != nil {
		h.logger.Printf("could not get vehicle: %v", err)
		http.Error(w, "could not get vehicle", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(vehicle); err != nil {
		h.logger.Printf("could not encode to json: %v", err)
		http.Error(w, "could not encode to json", http.StatusInternalServerError)
		return
	}
}

func (h VehicleHandler) Vehicles(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()

	vehicles, err := h.service.Vehicles(params.Get("model"), params.Get("brand"))
	if err != nil {
		h.logger.Printf("could not get vehicles: %v", err)
		http.Error(w, "could not get vehicles", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(vehicles); err != nil {
		h.logger.Printf("could not encode to json: %v", err)
		http.Error(w, "could not encode to json", http.StatusInternalServerError)
		return
	}
}

func (h VehicleHandler) CreateVehicle(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		h.logger.Printf("could not read from body: %v", err)
		http.Error(w, "could not read from body", http.StatusBadRequest)
		return
	}

	var v yakit.Vehicle
	if err := json.Unmarshal(body, &v); err != nil {
		h.logger.Printf("could not unmarshal json: %v", err)
		http.Error(w, "could not unmarshal json", http.StatusBadRequest)
		return
	}

	vehicle, err := h.service.CreateVehicle(v)
	if err != nil {
		h.logger.Printf("could not create vehicle: %v", err)
		http.Error(w, "could not create vehicle", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Location", fmt.Sprintf("/vehicles/%d", vehicle.ID))
	w.WriteHeader(http.StatusCreated)
}

func (h VehicleHandler) UpdateVehicle(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		h.logger.Printf("could not read from body: %v", err)
		http.Error(w, "could not read from body", http.StatusBadRequest)
		return
	}

	var v yakit.Vehicle
	if err := json.Unmarshal(body, &v); err != nil {
		h.logger.Printf("could not unmarshal json: %v", err)
		http.Error(w, "could not unmarshal json", http.StatusBadRequest)
		return
	}

	vars := mux.Vars(r)

	v.ID, err = strconv.Atoi(vars["id"])
	if err != nil {
		h.logger.Printf("could not convert ID %v to string: %v", v.ID, err)
		http.Error(w, "could not convert ID to string", http.StatusBadRequest)
		return
	}

	vehicle, err := h.service.UpdateVehicle(v)
	if err != nil {
		h.logger.Printf("could not update vehicle %d: %v", v.ID, err)
		http.Error(w, "could not update vehicle", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Location", fmt.Sprintf("/vehicles/%d", vehicle.ID))
	w.WriteHeader(http.StatusOK)
}

func (h VehicleHandler) DeleteVehicle(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	if err := h.service.DeleteVehicle(vars["id"]); err != nil {
		h.logger.Printf("could not delete vehicle %d: %v", vars["id"], err)
		http.Error(w, "could not delete", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
