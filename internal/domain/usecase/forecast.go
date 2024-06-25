package usecase

import (
	"context"
	"errors"

	"github.com/bruno-holanda15/weather_by_cep_goexpert/internal/domain/entity"
	"go.opentelemetry.io/otel/trace"
)

var (
	ErrorCanNotFindLocation         = errors.New("unable to find location by cep")
	ErrorExecutingRequestViaCep     = errors.New("error executing request to viacep")
	ErrorReadingBodyViaCep          = errors.New("error reading body from viacep response")
	ErrorUnsmarshalViaCep           = errors.New("error unmarshal from viacep body")
	ErrorExecutingRequestWeatherApi = errors.New("error executing request to weatherApi")
	ErrorReadingBodyWeatherApi      = errors.New("error reading body from weatherApi response")
	ErrorUnsmarshalWeatherApi       = errors.New("error unmarshal from weatherApi body")
	ErrorRemovingAccents            = errors.New("error removing accents")
)

type WeatherByCepUsecase struct {
	infosSearcher entity.InfosSearcherInterface
	tracer trace.Tracer
}

func NewWeatherByCepUsecase(infosSearcher entity.InfosSearcherInterface, trace trace.Tracer) *WeatherByCepUsecase {
	return &WeatherByCepUsecase{
		infosSearcher: infosSearcher,
		tracer: trace,
	}
}

func (w *WeatherByCepUsecase) Execute(ctx context.Context, input InputWbcUsecase) OutputWbcUsecase {
	var location entity.Location
	var cep string

	cep, err := w.infosSearcher.AddCep(input.Cep)
	if err != nil {
		return OutputWbcUsecase{Err: err}
	}
	location.Cep = cep

	var name string
	ctx, spanFindLocationName := w.tracer.Start(ctx, "find location name")
	name, err = w.infosSearcher.GetLocationName(cep)
	if err != nil {
		return OutputWbcUsecase{Err: err}
	}
	spanFindLocationName.End()
	location.Name = name

	var tempCelsius float32
	_, spanFindTemps := w.tracer.Start(ctx, "find temperatures")
	tempCelsius, err = w.infosSearcher.GetCelsiusTemp(name)
	if err != nil {
		return OutputWbcUsecase{Err: err}
	}
	spanFindTemps.End()
	location.TempCelsius = tempCelsius

	var tempFahrenheit, tempKelvin float32
	tempFahrenheit = w.infosSearcher.ConvertCelsiustoFahrenheit(tempCelsius)
	tempKelvin = w.infosSearcher.ConvertCelsiustoKelvin(tempCelsius)

	location.TempFahrenheit = tempFahrenheit
	location.TempKelvin = tempKelvin

	return OutputWbcUsecase{
		LocationName:   location.Name,
		TempCelsius:    location.TempCelsius,
		TempFahrenheit: location.TempFahrenheit,
		TempKelvin:     location.TempKelvin,
		Err:            nil,
	}
}
