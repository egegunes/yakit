package main

import (
	"log"
	"os"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/egegunes/yakit/database"
	"github.com/egegunes/yakit/server"
)

func main() {
	logger := log.New(os.Stderr, "yakit ", log.LstdFlags|log.Llongfile)

	logger.Printf("connecting to database Host:%s DB:%s User:%s", os.Getenv("DBHOST"), os.Getenv("DBNAME"), os.Getenv("DBUSER"))
	db, err := database.New(os.Getenv("DBHOST"), os.Getenv("DBNAME"), os.Getenv("DBUSER"), os.Getenv("DBPASS"))

	if err != nil {
		logger.Fatalf("connection to db failed: %v", err)
	}

	defer db.Close()

	r := mux.NewRouter()

	r.Handle("/metrics", promhttp.Handler())

	bs := database.NewBrandStore(db)
	bh := server.NewBrandHandler(bs, logger)

	r.HandleFunc("/brands", bh.Brands).Methods("GET")
	r.HandleFunc("/brands", bh.CreateBrand).Methods("POST")
	r.HandleFunc("/brands/{id:[0-9]+}", bh.Brand).Methods("GET")
	r.HandleFunc("/brands/{id:[0-9]+}", bh.UpdateBrand).Methods("POST")
	r.HandleFunc("/brands/{id:[0-9]+}", bh.DeleteBrand).Methods("DELETE")

	ms := database.NewModelStore(db)
	mh := server.NewModelHandler(ms, logger)

	r.HandleFunc("/models", mh.Models).Methods("GET")
	r.HandleFunc("/models", mh.CreateModel).Methods("POST")
	r.HandleFunc("/models/{id:[0-9]+}", mh.Model).Methods("GET")
	r.HandleFunc("/models/{id:[0-9]+}", mh.UpdateModel).Methods("POST")
	r.HandleFunc("/models/{id:[0-9]+}", mh.DeleteModel).Methods("DELETE")

	vs := database.NewVehicleStore(db)
	vh := server.NewVehicleHandler(vs, logger)

	r.HandleFunc("/vehicles", vh.Vehicles).Methods("GET")
	r.HandleFunc("/vehicles", vh.CreateVehicle).Methods("POST")
	r.HandleFunc("/vehicles/{id:[0-9]+}", vh.Vehicle).Methods("GET")
	r.HandleFunc("/vehicles/{id:[0-9]+}", vh.UpdateVehicle).Methods("POST")
	r.HandleFunc("/vehicles/{id:[0-9]+}", vh.DeleteVehicle).Methods("DELETE")

	listenAddr := os.Getenv("LISTENADDR")

	h1 := handlers.RecoveryHandler()(r)
	h2 := server.MetricsMiddleware(h1)
	h3 := handlers.LoggingHandler(os.Stdout, h2)
	h4 := handlers.CORS()(h3)

	srv := server.New(h4, listenAddr)

	logger.Printf("server starting on %s", listenAddr)
	err = srv.ListenAndServe()
	if err != nil {
		logger.Fatalf("server failed to start: %v", err)
	}
}
