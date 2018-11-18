package main

import (
	"log"
	"os"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"

	"yakit/database"
	"yakit/server"
)

func main() {
	logger := log.New(os.Stderr, "yakit ", log.LstdFlags|log.Llongfile)

	db := database.New(os.Getenv("DBHOST"), os.Getenv("DBNAME"), os.Getenv("DBUSER"), os.Getenv("DBPASS"))
	logger.Printf("connecting to database Host:%s DB:%s User:%s", os.Getenv("DBHOST"), os.Getenv("DBNAME"), os.Getenv("DBUSER"))
	conn, err := db.Open()

	if err != nil {
		logger.Fatalf("connection to db failed: %v", err)
	}

	defer conn.Close()

	r := mux.NewRouter()

	bs := database.BrandStore{DB: conn}
	bh := server.BrandHandler{Service: bs}

	r.HandleFunc("/brands", bh.Brands).Methods("GET")
	r.HandleFunc("/brands", bh.CreateBrand).Methods("POST")
	r.HandleFunc("/brands/{id:[0-9]+}", bh.Brand).Methods("GET")
	r.HandleFunc("/brands/{id:[0-9]+}", bh.UpdateBrand).Methods("POST")
	r.HandleFunc("/brands/{id:[0-9]+}", bh.DeleteBrand).Methods("DELETE")

	ms := database.NewModelStore(conn)
	mh := server.NewModelHandler(ms, logger)

	r.HandleFunc("/models", mh.Models).Methods("GET")
	r.HandleFunc("/models", mh.CreateModel).Methods("POST")
	r.HandleFunc("/models/{id:[0-9]+}", mh.Model).Methods("GET")
	r.HandleFunc("/models/{id:[0-9]+}", mh.UpdateModel).Methods("POST")
	r.HandleFunc("/models/{id:[0-9]+}", mh.DeleteModel).Methods("DELETE")

	vs := database.VehicleStore{DB: conn}
	vh := server.VehicleHandler{Service: vs}

	r.HandleFunc("/vehicles", vh.Vehicles).Methods("GET")
	r.HandleFunc("/vehicles", vh.CreateVehicle).Methods("POST")
	r.HandleFunc("/vehicles/{id:[0-9]+}", vh.Vehicle).Methods("GET")
	r.HandleFunc("/vehicles/{id:[0-9]+}", vh.UpdateVehicle).Methods("POST")
	r.HandleFunc("/vehicles/{id:[0-9]+}", vh.DeleteVehicle).Methods("DELETE")

	listenAddr := os.Getenv("LISTENADDR")

	h1 := handlers.LoggingHandler(os.Stdout, r)
	h2 := handlers.CORS()(h1)

	srv := server.New(h2, listenAddr)

	logger.Printf("server starting on %s", listenAddr)
	err = srv.ListenAndServe()
	if err != nil {
		logger.Fatalf("server failed to start: %v", err)
	}
}
