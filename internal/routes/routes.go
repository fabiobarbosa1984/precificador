package routes

import (
	"fmt"
	"net/http"
)

func NewRouter() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/", indexHandler)
	mux.HandleFunc("/api/data", apiDataHandler)

	return mux
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println(w, "API de precificação")
}

func apiDataHandler(w http.ResponseWriter, r *http.Request) {
	data := "dados da api"
	fmt.Println(w, data)
}
