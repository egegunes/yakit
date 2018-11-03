package main

import (
	"log"
	"os"

	"github.com/gorilla/mux"

	"yakit/database"
	"yakit/server"
)

func main() {
	logger := log.New(os.Stdout, "yakit ", log.LstdFlags|log.Llongfile)

	db := database.New(os.Getenv("DBHOST"), os.Getenv("DBNAME"), os.Getenv("DBUSER"), os.Getenv("DBPASS"))
	logger.Printf("connecting to database Host:%s DB:%s User:%s", os.Getenv("DBHOST"), os.Getenv("DBNAME"), os.Getenv("DBUSER"))
	conn, err := db.Open()

	if err != nil {
		log.Fatalf("connection to db failed: %v", err)
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

	ms := database.ModelStore{DB: conn}
	mh := server.ModelHandler{Service: ms}

	r.HandleFunc("/models", mh.Models).Methods("GET")
	r.HandleFunc("/models", mh.CreateModel).Methods("POST")
	r.HandleFunc("/models/{id:[0-9]+}", mh.Model).Methods("GET")
	r.HandleFunc("/models/{id:[0-9]+}", mh.UpdateModel).Methods("POST")

	listenAddr := os.Getenv("LISTENADDR")

	srv := server.New(r, listenAddr)

	logger.Printf("server starting on %s", listenAddr)
	err = srv.ListenAndServe()
	if err != nil {
		log.Fatalf("server failed to start: %v", err)
	}
}
