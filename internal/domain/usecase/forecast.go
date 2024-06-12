package usecase

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"unicode"

	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

type InputWbcUsecase struct {
	Cep string
}

type OutputWbcUsecase struct {
	Localidade  string
	Temperatura float32
}

type ViaCepInfo struct {
	Location string `json:"localidade"`
}

type WeatherApiInfo struct {
	Current Current `json:"current"`
}

type Current struct {
	TempCelsius float32 `json:"temp_c"`
}

type WeatherByCepUsecase struct{}

func (w *WeatherByCepUsecase) Execute(input InputWbcUsecase) (OutputWbcUsecase, error) {

	location, err := getLocation(input.Cep)
	if err != nil {
		return OutputWbcUsecase{}, err
	}

	celsiusTemp, err := getCelsiusTemp(location)
	if err != nil {
		return OutputWbcUsecase{}, err
	}

	return OutputWbcUsecase{
		Localidade:  location,
		Temperatura: celsiusTemp,
	}, nil
}

func getLocation(cep string) (string, error) {
	resViaCep, err := http.Get("http://viacep.com.br/ws/" + cep + "/json/")

	if err != nil {
		return "", err
	}
	defer resViaCep.Body.Close()

	body, err := io.ReadAll(resViaCep.Body)
	if err != nil {
		return "", err
	}

	var viaCepInfo ViaCepInfo
	err = json.Unmarshal(body, &viaCepInfo)
	if err != nil {
		return "", err
	}
	fmt.Println(body, viaCepInfo, "aqui")

	viaCepInfo.Location, err = removeAccents(viaCepInfo.Location)
	if err != nil {
		return "", err
	}

	return viaCepInfo.Location, nil
}

func getCelsiusTemp(location string) (float32, error) {
	urlWeather := fmt.Sprintf("http://api.weatherapi.com/v1/current.json?key=ad23180c91fc43eba2d232304241106&q=%s", strings.Replace(location, " ", "+", -1))

	fmt.Println(urlWeather)
	resWeatherApi, err := http.Get(urlWeather)
	if err != nil {
		return 0, err
	}
	defer resWeatherApi.Body.Close()

	bodyWeather, err := io.ReadAll(resWeatherApi.Body)
	fmt.Println(string(bodyWeather))
	if err != nil {
		return 0, err
	}

	var weatherInfo WeatherApiInfo
	err = json.Unmarshal(bodyWeather, &weatherInfo)
	if err != nil {
		return 0, err
	}

	return weatherInfo.Current.TempCelsius, nil
}

var normalizer = transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)

func removeAccents(str string) (string, error) {
	s, _, err := transform.String(normalizer, str)
	if err != nil {
		return "", err
	}
	return s, err
}
