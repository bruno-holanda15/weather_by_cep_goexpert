package usecase

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type InputWbcUsecase struct {
	Cep string
}

type OutputWbcUsecase struct {
	Temperatura   float32
	ResViaCep     any
	ResWeatherApi any
}

type ViaCepInfos struct {
	Localidade string
}

type WeatherByCepUsecase struct{}

func (w *WeatherByCepUsecase) Execute(input InputWbcUsecase) (OutputWbcUsecase, error) {

	resViaCep, err := http.Get("http://viacep.com.br/ws/" + input.Cep + "/json/")
	if err != nil {
		return OutputWbcUsecase{}, err
	}
	defer resViaCep.Body.Close()

	body, err := io.ReadAll(resViaCep.Body)
	if err != nil {
		return OutputWbcUsecase{}, err
	}

	var location ViaCepInfos
	err = json.Unmarshal(body, &location)
	if err != nil {
		return OutputWbcUsecase{}, err
	}

	urlWeather := fmt.Sprintf("http://api.weatherapi.com/v1/current.json?key=ad23180c91fc43eba2d232304241106&q=%s", strings.Replace(location.Localidade, " ", "+",-1))

	fmt.Println(urlWeather)
	resWeatherApi, err := http.Get(urlWeather)
	if err != nil {
		return OutputWbcUsecase{}, err
	}
	defer resWeatherApi.Body.Close()

	bodyWeather, err := io.ReadAll(resWeatherApi.Body)
	if err != nil {
		return OutputWbcUsecase{}, err
	}

	return OutputWbcUsecase{
		Temperatura: 12.2,
		ResViaCep:   string(body),
		ResWeatherApi: string(bodyWeather),
	}, nil
}
