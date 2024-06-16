package entity

type Location struct {
	Cep            string
	Name           string
	TempCelsius    float32
	TempFahrenheit float32
	TempKelvin     float32
}

type InfosSearcher struct {}

type ViaCepInfo struct {
	LocationName string `json:"localidade"`
}

type Current struct {
	TempCelsius float32 `json:"temp_c"`
}

type WeatherApiInfo struct {
	Current Current `json:"current"`
}
