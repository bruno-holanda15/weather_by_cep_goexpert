package usecase

type InputWbcUsecase struct {
	Cep string
}

type OutputWbcUsecase struct {
	TempCelsius    float32 `json:"temp_C"`
	TempFahrenheit float32 `json:"temp_F"`
	TempKelvin     float32 `json:"temp_K"`
	Err            error
}
