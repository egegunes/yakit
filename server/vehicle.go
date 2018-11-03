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

type VehicleHandler struct {
	Service yakit.VehicleService
}

func (h VehicleHandler) Vehicle(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	vehicle, err := h.Service.Vehicle(vars["id"])
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error: %s", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(vehicle)
}

func (h VehicleHandler) Vehicles(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()

	vehicles, err := h.Service.Vehicles(params.Get("model"), params.Get("brand"))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error: %s", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(vehicles)
}

func (h VehicleHandler) CreateVehicle(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Error: %s", err)
		return
	}

	var v yakit.Vehicle
	err = json.Unmarshal(body, &v)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Error: %s", err)
		return
	}

	vehicle, err := h.Service.CreateVehicle(v)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error: %s", err)
		return
	}

	w.Header().Set("Location", fmt.Sprintf("/vehicles/%d", vehicle.ID))
	w.WriteHeader(http.StatusCreated)
}

func (h VehicleHandler) UpdateVehicle(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Error: %s", err)
		return
	}

	var v yakit.Vehicle
	err = json.Unmarshal(body, &v)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Error: %s", err)
		return
	}

	vars := mux.Vars(r)

	v.ID, err = strconv.Atoi(vars["id"])

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Error: %s", err)
		return
	}

	vehicle, err := h.Service.UpdateVehicle(v)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error: %s", err)
		return
	}

	w.Header().Set("Location", fmt.Sprintf("/vehicles/%d", vehicle.ID))
	w.WriteHeader(http.StatusOK)
}

func (h VehicleHandler) DeleteVehicle(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	err := h.Service.DeleteVehicle(vars["id"])
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error: %s", err)
		return
	}

	w.WriteHeader(http.StatusOK)
}
