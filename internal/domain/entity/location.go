package entity

import (
	"errors"
	"strings"
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
		return errors.New("invalid cep")
	}

	return nil
}
