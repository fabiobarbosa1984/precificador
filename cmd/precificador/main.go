package main

import (
	"log"
	"time"

	"github.com/fabiobarbosa1984/precificador/internal/calculadora"
)

func main() {
	liquidacao := time.Date(2023, time.August, 30, 0, 0, 0, 0, time.UTC)
	vencimento := time.Date(2025, time.July, 1, 0, 0, 0, 0, time.UTC)

	novoCalculo := calculadora.NovoCalculo(calculadora.NTN_F, vencimento, liquidacao)

	//	novoCalculo.PrecificarLTN()

	//log.Println(novoCalculo.Preco)
	log.Println(novoCalculo.Titulo.Cupom)

	/* router := routes.NewRouter()

	port := 8080
	addr := fmt.Sprintf(":%d", port)
	log.Println("Servidor rodando em http://localhost", addr)

	err := http.ListenAndServe(addr, router)
	if err != nil {
		panic(err)
	} */
}
