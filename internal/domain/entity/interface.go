package entity

type InfosSearcherInterface interface {
	AddCep(string) (string, error)
	ConvertCelsiustoFahrenheit(float32) float32
	ConvertCelsiustoKelvin(float32) float32
	GetLocationName(string) (string, error)
	GetCelsiusTemp(string) (float32, error)
}
