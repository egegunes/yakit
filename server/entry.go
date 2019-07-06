package server

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strconv"

	// "fmt"
	// "io/ioutil"
	"log"
	"net/http"

	// "strconv"

	"github.com/gorilla/mux"

	"github.com/egegunes/yakit/yakit"
)

type EntryHandler struct {
	logger  *log.Logger
	service yakit.EntryService
}

func NewEntryHandler(s yakit.EntryService, l *log.Logger) *EntryHandler {
	return &EntryHandler{service: s, logger: l}
}

func (h EntryHandler) Entry(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	entry, err := h.service.Entry(vars["id"])
	if err != nil {
		h.logger.Printf("could not get entry: %v", err)
		http.Error(w, "could not get entry", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(entry); err != nil {
		h.logger.Printf("could not encode to json: %v", err)
		http.Error(w, "could not encode to json", http.StatusInternalServerError)
		return
	}
}

func (h EntryHandler) Entries(w http.ResponseWriter, r *http.Request) {
	entries, err := h.service.Entries()
	if err != nil {
		h.logger.Printf("could not get entries: %v", err)
		http.Error(w, "could not get entries", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(entries); err != nil {
		h.logger.Printf("could not encode to json: %v", err)
		http.Error(w, "could not encode to json", http.StatusInternalServerError)
		return
	}
}

func (h EntryHandler) CreateEntry(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		h.logger.Printf("could not read from body: %v", err)
		http.Error(w, "could not read from body", http.StatusBadRequest)
		return
	}

	var e yakit.Entry
	if err := json.Unmarshal(body, &e); err != nil {
		h.logger.Printf("could not unmarshal json: %v", err)
		http.Error(w, "could not unmarshal json", http.StatusBadRequest)
		return
	}

	entry, err := h.service.CreateEntry(e)
	if err != nil {
		h.logger.Printf("could not create entry: %v", err)
		http.Error(w, "could not create entry", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Location", fmt.Sprintf("/entries/%d", entry.ID))
	w.WriteHeader(http.StatusCreated)
}

func (h EntryHandler) UpdateEntry(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		h.logger.Printf("could not read from body: %v", err)
		http.Error(w, "could not read from body", http.StatusBadRequest)
		return
	}

	var e yakit.Entry
	if err := json.Unmarshal(body, &e); err != nil {
		h.logger.Printf("could not unmarshal json: %v", err)
		http.Error(w, "could not unmarshal json", http.StatusBadRequest)
		return
	}

	vars := mux.Vars(r)

	e.ID, err = strconv.Atoi(vars["id"])
	if err != nil {
		h.logger.Printf("could not convert ID %v to string: %v", e.ID, err)
		http.Error(w, "could not convert ID to string", http.StatusBadRequest)
		return
	}

	entry, err := h.service.UpdateEntry(e)
	if err != nil {
		h.logger.Printf("could not update entry %d: %v", e.ID, err)
		http.Error(w, "could not update entry", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Location", fmt.Sprintf("/entries/%d", entry.ID))
	w.WriteHeader(http.StatusOK)
}

func (h EntryHandler) DeleteEntry(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	if err := h.service.DeleteEntry(vars["id"]); err != nil {
		h.logger.Printf("could not delete entry %d: %v", vars["id"], err)
		http.Error(w, "could not delete", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
