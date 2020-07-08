package main

import (
	"github.com/anirvansen/golang_rest_api/handlers"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {

	dbDriver := "mysql"
	dbUser := os.Getenv("dbUser")
	dbPass := os.Getenv("dbPass")
	dbName := os.Getenv("dbName")

	dbProperties := map[string]string{
		"dbDriver" : dbDriver,
		"dbUser" : dbUser,
		"dbPass" : dbPass,
		"dbName" : dbName,
	}


	logger :=  log.New(os.Stdout,"golang-rest-api",log.LstdFlags)


	sm := mux.NewRouter()

	ph := handlers.ProductHandler(logger,dbProperties)

	getRouter := sm.Methods(http.MethodGet).Subrouter()
	getRouter.HandleFunc("/products",ph.GetProducts)
	getRouter.HandleFunc("/products/{id:[0-9]+}",ph.GetProductById)

	postRouter := sm.Methods(http.MethodPost).Subrouter()
	postRouter.HandleFunc("/products",ph.SaveProduct)

	putRouter := sm.Methods(http.MethodPut).Subrouter()
	putRouter.HandleFunc("/products/{id:[0-9]+}",ph.UpdateProductById)

	deleteRouter := sm.Methods(http.MethodDelete).Subrouter()
	deleteRouter.HandleFunc("/products/{id:[0-9]+}",ph.DeleteProductById)




	s := &http.Server{
		Addr:              ":9000",
		Handler:           sm,
		ReadTimeout:       1 * time.Second,
		WriteTimeout:      1 * time.Second,
		IdleTimeout:       120* time.Second,
	}

	s.ListenAndServe()

}
