package usecase

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"unicode"

	"github.com/bruno-holanda15/weather_by_cep_goexpert/configs"
	"github.com/bruno-holanda15/weather_by_cep_goexpert/internal/domain/entity"
	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

var (
	CanNotFindLocation = errors.New("unable to find location by cep")
	ErrorExecutingRequestViaCep = errors.New("error executing request to viacep")
	ErrorReadingBodyViaCep = errors.New("error reading body from viacep response")
	ErrorUnsmarshalViaCep = errors.New("error unmarshal from viacep body")
	ErrorExecutingRequestWeatherApi = errors.New("error executing request to weatherApi")
	ErrorReadingBodyWeatherApi = errors.New("error reading body from weatherApi response")
	ErrorUnsmarshalWeatherApi = errors.New("error unmarshal from weatherApi body")
	ErrorRemovingAccents = errors.New("error removing accents")
)

type InputWbcUsecase struct {
	Cep string
}

type OutputWbcUsecase struct {
	TempCelsius    float32 `json:"temp_C"`
	TempFahrenheit float32 `json:"temp_F"`
	TempKelvin     float32 `json:"temp_K"`
}

type ViaCepInfo struct {
	LocationName string `json:"localidade"`
}

type WeatherApiInfo struct {
	Current Current `json:"current"`
}

type Current struct {
	TempCelsius float32 `json:"temp_c"`
}

type WeatherByCepUsecase struct{}

func (w *WeatherByCepUsecase) Execute(input InputWbcUsecase) (OutputWbcUsecase, error) {
	location := entity.NewLocation()
	err := location.AddCep(input.Cep)
	if err != nil {
		return OutputWbcUsecase{}, err
	}

	location.Name, err = getLocationName(input.Cep)
	if err != nil {
		return OutputWbcUsecase{}, err
	}

	location.TempCelsius, err = getCelsiusTemp(location.Name)
	location.FillOtherTempsFromCelsius()
	if err != nil {
		return OutputWbcUsecase{}, err
	}

	return OutputWbcUsecase{
		TempCelsius: location.TempCelsius,
		TempFahrenheit: location.TempFahrenheit,
		TempKelvin: location.TempKelvin,
	}, nil
}

func getLocationName(cep string) (string, error) {
	resViaCep, err := http.Get("http://viacep.com.br/ws/" + cep + "/json/")

	if err != nil {
		return "", ErrorExecutingRequestViaCep
	}
	defer resViaCep.Body.Close()

	body, err := io.ReadAll(resViaCep.Body)
	if err != nil {
		return "", ErrorReadingBodyViaCep
	}

	var viaCepInfo ViaCepInfo
	err = json.Unmarshal(body, &viaCepInfo)
	if err != nil {
		return "", ErrorUnsmarshalViaCep
	}

	if viaCepInfo.LocationName == "" {
		return "", CanNotFindLocation
	}

	viaCepInfo.LocationName, err = removeAccents(viaCepInfo.LocationName)
	if err != nil {
		return "", ErrorRemovingAccents
	}

	return viaCepInfo.LocationName, nil
}

func getCelsiusTemp(location string) (float32, error) {
	apiToken := configs.Env("WEATHER_TOKEN")
	urlWeather := fmt.Sprintf("http://api.weatherapi.com/v1/current.json?key="+apiToken+"&q=%s", strings.Replace(location, " ", "+", -1))

	resWeatherApi, err := http.Get(urlWeather)
	if err != nil {
		return 0, ErrorExecutingRequestWeatherApi
	}
	defer resWeatherApi.Body.Close()

	bodyWeather, err := io.ReadAll(resWeatherApi.Body)
	if err != nil {
		return 0, ErrorReadingBodyWeatherApi
	}

	var weatherInfo WeatherApiInfo
	err = json.Unmarshal(bodyWeather, &weatherInfo)
	if err != nil {
		return 0, ErrorUnsmarshalWeatherApi
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
