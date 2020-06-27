package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// verificando se o usuario informou uma URL válida
func notFound(w http.ResponseWriter, r *http.Request) {
	var infoErro mensagemErro
	if 404 == http.StatusNotFound {
		infoErro.Status = "Erro: " + http.StatusText(http.StatusNotFound)
		infoErro.Informacao = "Verifique se a URL foi informada corretamente"
		infoErro.Codigo = http.StatusNotFound

		dataErro, _ := json.Marshal(infoErro)

		w.Header().Set("Content-Type", "application/json")
		w.Write(dataErro)
		return
	}
}
func main() {
	router := mux.NewRouter()
	router.HandleFunc("/cnpj/{id}", GetCnpjEndpoint).Methods("GET") //URL: localhost:8000/cnpj/17184406000174
	router.HandleFunc("/quote", PostFreteEndpoint).Methods("POST")  //URL: localhost:8000/quote
	router.NotFoundHandler = http.HandlerFunc(notFound)
	log.Fatal(http.ListenAndServe(":8000", router))
}
