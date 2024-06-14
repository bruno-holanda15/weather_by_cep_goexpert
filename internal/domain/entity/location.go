package entity

import (
	"errors"
	"strings"
)

var (
	InvalidCep = errors.New("invalid cep")
	ErrorCelsiusEmpty = errors.New("temp Celsius not completed yet")
)

type Location struct {
	Cep            string
	Name           string
	TempCelsius    float32
	TempFahrenheit float32
	TempKelvin     float32
}

func NewLocation() *Location {
	return &Location{}
}

func (l *Location) AddCep(cep string) error {
	if err := isCepValid(cep); err != nil {
		return err
	}
	l.Cep = cep
	return nil
}

func isCepValid(cep string) error {
	if strings.NewReader(cep).Size() != 8 {
		return InvalidCep
	}

	return nil
}

func (l *Location) FillOtherTempsFromCelsius() error {
	if l.TempCelsius == 0 {
		return ErrorCelsiusEmpty
	}

	l.TempFahrenheit = l.TempCelsius*1.8 + 32
	l.TempKelvin = l.TempCelsius + 273
	return nil
}
