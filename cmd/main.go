package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/bruno-holanda15/weather_by_cep_goexpert/configs"
	"github.com/bruno-holanda15/weather_by_cep_goexpert/internal/domain/usecase"
	"github.com/bruno-holanda15/weather_by_cep_goexpert/internal/infra/web"
)

func main() {
	weatherByCEPUsecase := usecase.WeatherByCepUsecase{}
	wbcHandler := web.NewWeatherByCepHttp(weatherByCEPUsecase)

	http.Handle("/weather/{cep}", wbcHandler)

	fmt.Println("Listening http server http://localhost:8083")
	if err := http.ListenAndServe(":8083", nil); err != nil {
		log.Fatalf("error starting http server - %v", err)
	}
}

func init() {
	loader := &configs.Loader{}
	loader.LoadEnv()
}