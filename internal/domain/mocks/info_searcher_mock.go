package mocks

import "errors"

type InfoSearcherMock struct {
	AddCepFunc                     func(string) (string, error)
	GetLocationNameFunc            func(string) (string, error)
	GetCelsiusTempFunc             func(string) (float32, error)
	ConvertCelsiustoFahrenheitFunc func(float32) float32
	ConvertCelsiustoKelvinFunc     func(float32) float32
}

func NewInfoSearcherMock() *InfoSearcherMock {
	return &InfoSearcherMock{
		AddCepFunc:                     func(s string) (string, error) { return "", nil },
		GetLocationNameFunc:            func(s string) (string, error) { return "", nil },
		GetCelsiusTempFunc:             func(s string) (float32, error) { return 0, nil },
		ConvertCelsiustoFahrenheitFunc: func(f float32) float32 { return 0 },
		ConvertCelsiustoKelvinFunc:     func(f float32) float32 { return 0 },
	}
}

func (i *InfoSearcherMock) AddCep(cep string) (string, error) {
	if i.AddCepFunc != nil {
		return i.AddCepFunc(cep)
	}
	return "", errors.New("func not implemented")
}

func (i *InfoSearcherMock) GetLocationName(cep string) (string, error) {
	if i.GetLocationNameFunc != nil {
		return i.GetLocationNameFunc(cep)
	}
	return "", errors.New("func not implemented")
}

func (i *InfoSearcherMock) GetCelsiusTemp(name string) (float32, error) {
	if i.GetCelsiusTempFunc != nil {
		return i.GetCelsiusTempFunc(name)
	}
	return 0, errors.New("func not implemented")
}

func (i *InfoSearcherMock) ConvertCelsiustoFahrenheit(celsius float32) float32 {
	if i.ConvertCelsiustoFahrenheitFunc != nil {
		return i.ConvertCelsiustoFahrenheitFunc(celsius)
	}
	return 0
}

func (i *InfoSearcherMock) ConvertCelsiustoKelvin(celsius float32) float32 {
	if i.ConvertCelsiustoKelvinFunc != nil {
		return i.ConvertCelsiustoKelvinFunc(celsius)
	}

	return 0
}
