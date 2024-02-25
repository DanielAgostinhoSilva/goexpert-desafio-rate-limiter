package main

import (
	"fmt"
	"github.com/DanielAgostinhoSilva/goexpert-desafio-rate-limiter/src/infrastructure/env"
	"net/http"
)

func helloHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello, World!")
}

func main() {
	env.LoadConfig("./.env")

	//http.HandleFunc("/", helloHandler)
	//
	//fmt.Println("Server is running on port 8080")
	//if err := http.ListenAndServe(":8080", nil); err != nil {
	//	panic(err)
	//}
}
