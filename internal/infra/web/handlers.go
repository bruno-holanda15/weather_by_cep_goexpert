package web

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/bruno-holanda15/weather_by_cep_goexpert/internal/domain/entity"
	"github.com/bruno-holanda15/weather_by_cep_goexpert/internal/domain/usecase"
)

type WeatherByCepHttp struct {
	usecase *usecase.WeatherByCepUsecase
}

func NewWeatherByCepHttp(usecase *usecase.WeatherByCepUsecase) *WeatherByCepHttp {
	return &WeatherByCepHttp{
		usecase: usecase,
	}
}

func (we *WeatherByCepHttp) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	cep := r.PathValue("cep")
	input := usecase.InputWbcUsecase{Cep: cep}

	output := we.usecase.Execute(input)
	if output.Err != nil {
		err := output.Err
		if err == entity.ErrorCanNotFindLocation ||  err == entity.ErrorInvalidCep {
			w.WriteHeader(http.StatusUnprocessableEntity)
			w.Write([]byte(err.Error()))
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Println(output)
	json.NewEncoder(w).Encode(output)
}

type CepBody struct {
	Cep string `json:"cep"`
}

func ValidateCep(w http.ResponseWriter, r *http.Request) {
	var reqBody CepBody

	err := json.NewDecoder(r.Body).Decode(&reqBody)
	if err != nil || reqBody.Cep == "" {
		http.Error(w, "Invalid request body", http.StatusUnprocessableEntity)
		return
	}

	err = entity.IsCepValid(reqBody.Cep)
	if err != nil {
		http.Error(w, "Invalid cep", http.StatusUnprocessableEntity)
		return
	}

	reqWeather, err := http.NewRequest("GET", "http://goapp2:8082/weather/"+reqBody.Cep, nil)
	if err != nil {
		http.Error(w, "error preparing request to go_wbc2", http.StatusInternalServerError)
		return
	}

	resWeather, err := http.DefaultClient.Do(reqWeather)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resWeather.Body.Close()

	if resWeather.StatusCode != 200 {
		w.WriteHeader(resWeather.StatusCode)
		return
	}

	weatherBody, err := io.ReadAll(resWeather.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(weatherBody)
}
