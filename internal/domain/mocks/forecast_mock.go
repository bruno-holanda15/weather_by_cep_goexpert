package mocks

type WeatherByCepUsecaseMock struct {
	GetLocationNameFunc func(cep string) (string, error)
	GetCelsiusTempFunc  func(location string) (float32, error)
}

func NewWeatherByCepUsecaseMock() *WeatherByCepUsecaseMock {
	return &WeatherByCepUsecaseMock{
		GetLocationNameFunc: func(cep string) (string, error) { return "", nil },
		GetCelsiusTempFunc:  func(location string) (float32, error) { return 0, nil },
	}
}

func (w *WeatherByCepUsecaseMock) GetLocationName(cep string) (string, error) {
	return w.GetLocationNameFunc(cep)
}

func (w *WeatherByCepUsecaseMock) GetCelsiusTemp(location string) (float32, error) {
	return w.GetCelsiusTempFunc(location)
}
