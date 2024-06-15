package usecase

import (
	"testing"

	"github.com/bruno-holanda15/weather_by_cep_goexpert/internal/domain/mocks"
	"github.com/stretchr/testify/assert"
)

func TestWeatherByCepUsecase_Execu(t *testing.T) {
	input := InputWbcUsecase{Cep: "88010290"}
	// inputInvalid := InputWbcUsecase{Cep: "92019201202"}

	tests := []struct {
		name string
		wbcUseCaseMock func(*mocks.WeatherByCepUsecaseMock)
		expected *OutputWbcUsecase
	} {
		{
			name: "Sucessful Search", 
			wbcUseCaseMock: func(wmock *mocks.WeatherByCepUsecaseMock) {
				wmock.GetLocationNameFunc = func(cep string) (string, error) {
					return "Florianopolis", nil
				}
				wmock.GetCelsiusTempFunc = func(location string) (float32, error) {
					return 20.0, nil
				}
			},
			expected: &OutputWbcUsecase{
				TempCelsius: 20.0,
				TempFahrenheit: 68.0,
				TempKelvin: 293.0,
				Err: nil,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.wbcUseCaseMock(mocks.NewWeatherByCepUsecaseMock())

			usecase := WeatherByCepUsecase{}
			output := usecase.Execute(input)
			assert.Equal(t, tt.expected, output)

		})
	}
}
