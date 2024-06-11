package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/bruno-holanda15/weather_by_cep_goexpert/internal/domain/usecase"
)

func main() {

	weatherByCEPUsecase := usecase.WeatherByCepUsecase{}

	http.HandleFunc("/weather/{cep}", func(w http.ResponseWriter, r *http.Request) {
		cep := r.PathValue("cep")
		input := usecase.InputWbcUsecase{Cep: cep}

		output, err := weatherByCEPUsecase.Execute(input)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		fmt.Println(output)
		w.Write([]byte("Requisição viaCep sucesso\n"))
	})

	fmt.Println("Listening http server http://localhost:8083")
	if err := http.ListenAndServe(":8083", nil); err != nil {
		log.Fatalf("error starting http server - %v", err)
	}
}
