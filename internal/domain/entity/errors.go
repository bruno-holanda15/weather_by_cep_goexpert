package entity

import "errors"

var (
	ErrorInvalidCep                 = errors.New("invalid cep")
	ErrorCelsiusEmpty               = errors.New("temp Celsius not completed yet")
	ErrorCanNotFindLocation         = errors.New("unable to find location by cep")
	ErrorExecutingRequestViaCep     = errors.New("error executing request to viacep")
	ErrorReadingBodyViaCep          = errors.New("error reading body from viacep response")
	ErrorUnsmarshalViaCep           = errors.New("error unmarshal from viacep body")
	ErrorExecutingRequestWeatherApi = errors.New("error executing request to weatherApi")
	ErrorReadingBodyWeatherApi      = errors.New("error reading body from weatherApi response")
	ErrorUnsmarshalWeatherApi       = errors.New("error unmarshal from weatherApi body")
	ErrorRemovingAccents            = errors.New("error removing accents")
	ErrorEmptyLocationName          = errors.New("empty location name")
	ErrorEmptyCep                   = errors.New("empty cep")
)
