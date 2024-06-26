package usecase

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/bruno-holanda15/weather_by_cep_goexpert/internal/domain/entity"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
)

type ValidateCepUsecase struct{}

func (v *ValidateCepUsecase) Execute(ctx context.Context, input InputWbcUsecase) OutputWbcUsecase {

	err := entity.IsCepValid(input.Cep)
	if err != nil {
		return OutputWbcUsecase{Err: err}
	}

	reqWeather, err := http.NewRequestWithContext(ctx, "GET", "http://goapp2:8082/weather/"+input.Cep, nil)
	if err != nil {
		return OutputWbcUsecase{Err: err}
	}

	otel.GetTextMapPropagator().Inject(ctx, propagation.HeaderCarrier(reqWeather.Header))
	resWeather, err := http.DefaultClient.Do(reqWeather)
	if err != nil {
		return OutputWbcUsecase{Err: err}
	}
	defer resWeather.Body.Close()

	if resWeather.StatusCode == 404 {
		return OutputWbcUsecase{Err: entity.ErrorEmptyCep}
	}

	var output OutputWbcUsecase
	json.NewDecoder(resWeather.Body).Decode(&output)

	return output
}
