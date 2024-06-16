package entity

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"unicode"

	"github.com/bruno-holanda15/weather_by_cep_goexpert/configs"
	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

func NewInfosSearcher() *InfosSearcher {
	return &InfosSearcher{}
}

func (l *InfosSearcher) AddCep(cep string) (string, error) {
	err := isCepValid(cep)
	if err != nil {
		return "", err
	}

	return cep, err
}

func isCepValid(cep string) error {
	if strings.NewReader(cep).Size() != 8 {
		return ErrorInvalidCep
	}

	if _, err := strconv.Atoi(cep); err != nil {
		return ErrorInvalidCep
	}

	return nil
}

func (l *InfosSearcher) GetLocationName(cep string) (string, error) {
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
		return "", ErrorCanNotFindLocation
	}

	viaCepInfo.LocationName, err = removeAccents(viaCepInfo.LocationName)
	if err != nil {
		return "", ErrorRemovingAccents
	}

	return viaCepInfo.LocationName, nil
}

func (l *InfosSearcher) GetCelsiusTemp(locationName string) (float32, error) {
	apiToken := configs.Env("WEATHER_TOKEN", "Teste")
	urlWeather := fmt.Sprintf("http://api.weatherapi.com/v1/current.json?key="+apiToken+"&q=%s", strings.Replace(locationName, " ", "+", -1))

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

func (l *InfosSearcher) ConvertCelsiustoFahrenheit(c float32) float32 {
	return c *1.8 + 32
}

func (l *InfosSearcher) ConvertCelsiustoKelvin(c float32) float32 {
	return c  + 273
}
