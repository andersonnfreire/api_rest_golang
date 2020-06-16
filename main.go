package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/cnpj/{id}", GetCnpjEndpoint).Methods("GET") //URL: localhost:8000/cnpj/27865757000102
	router.HandleFunc("/quote", PostFreteEndpoint).Methods("POST")  //URL: localhost:8000/quote
	log.Fatal(http.ListenAndServe(":8000", router))
}
