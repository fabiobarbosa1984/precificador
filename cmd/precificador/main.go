package main

import (
	"fmt"
	"log"
	"net/http"
	"precificador/internal/routes"
)

func main() {
	router := routes.NewRouter()

	port := 8080
	addr := fmt.Sprintf(":%d", port)
	log.Println("Servidor rodando em http://localhost", addr)

	err := http.ListenAndServe(addr, router)
	if err != nil {
		panic(err)
	}
}
