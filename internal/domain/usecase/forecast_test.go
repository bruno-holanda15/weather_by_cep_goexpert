package usecase

import (
	"testing"

	"github.com/bruno-holanda15/weather_by_cep_goexpert/internal/domain/entity"
	"github.com/bruno-holanda15/weather_by_cep_goexpert/internal/domain/mocks"
	"github.com/stretchr/testify/assert"
)

func TestWeatherByCepUsecase_Execute(t *testing.T) {
	inputValid := InputWbcUsecase{Cep: "88010290"}
	inputInvalid := InputWbcUsecase{Cep: "toptop"}
	inputNotFound := InputWbcUsecase{Cep: "01234567"}

	tests := []struct {
		name     string
		searcher func(*mocks.InfoSearcherMock)
		expected OutputWbcUsecase
	}{
		{
			name: "Valid Cep",
			searcher: func(searcher *mocks.InfoSearcherMock) {
				searcher.AddCepFunc = func(cep string) (string, error) {
					return "88010290", nil
				}
				searcher.GetLocationNameFunc = func(cep string) (string, error) {
					return "Florianópolis", nil
				}
				searcher.GetCelsiusTempFunc = func(name string) (float32, error) {
					return 20, nil
				}
				searcher.ConvertCelsiustoFahrenheitFunc = func(f float32) float32 {
					return 68
				}
				searcher.ConvertCelsiustoKelvinFunc = func(f float32) float32 {
					return 293
				}
			},
			expected: OutputWbcUsecase{
				TempCelsius:    20,
				TempFahrenheit: 68,
				TempKelvin:     293,
				Err:            nil,
			},
		},
		{
			name: "Invalid Cep",
			searcher: func(searcher *mocks.InfoSearcherMock) {
				searcher.AddCepFunc = func(cep string) (string, error) {
					return "", entity.ErrorInvalidCep
				}
			},
			expected: OutputWbcUsecase{
				Err: entity.ErrorInvalidCep,
			},
		},
		{
			name: "Cep not found",
			searcher: func(searcher *mocks.InfoSearcherMock) {
				searcher.AddCepFunc = func(cep string) (string, error) {
					return "01234567", nil
				}
				searcher.GetLocationNameFunc = func(cep string) (string, error) {
					return "", entity.ErrorCanNotFindLocation
				}
			},
			expected: OutputWbcUsecase{
				Err: entity.ErrorCanNotFindLocation,
			},
		},
		{
			name: "Error request to weather api",
			searcher: func(searcher *mocks.InfoSearcherMock) {
				searcher.AddCepFunc = func(cep string) (string, error) {
					return "88010290", nil
				}
				searcher.GetLocationNameFunc = func(cep string) (string, error) {
					return "Florianópolis", nil
				}
				searcher.GetCelsiusTempFunc = func(name string) (float32, error) {
					return 0, entity.ErrorExecutingRequestWeatherApi
				}
			},
			expected: OutputWbcUsecase{
				Err: entity.ErrorExecutingRequestWeatherApi,
			},
		},
	}

	for _, tt := range tests {
		mock := mocks.NewInfoSearcherMock()
		tt.searcher(mock)
		usecase := NewWeatherByCepUsecase(mock)

		if tt.name == "Valid Cep" {
			output := usecase.Execute(inputValid)
			assert.Equal(t, tt.expected, output)
		} else if tt.name == "Invalid Cep" {
			output := usecase.Execute(inputInvalid)
			assert.Equal(t, tt.expected, output)
		} else if tt.name == "Cep not found" {
			output := usecase.Execute(inputNotFound)
			assert.Equal(t, tt.expected, output)
		} else {
			output := usecase.Execute(inputValid)
			assert.Equal(t, tt.expected, output)
		}
	}
}
