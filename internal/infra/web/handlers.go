package web

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/bruno-holanda15/weather_by_cep_goexpert/internal/domain/entity"
	"github.com/bruno-holanda15/weather_by_cep_goexpert/internal/domain/usecase"
)

type WeatherByCepHttp struct {
	usecase usecase.WeatherByCepUsecase
}

func NewWeatherByCepHttp(usecase usecase.WeatherByCepUsecase) *WeatherByCepHttp {
	return &WeatherByCepHttp{
		usecase: usecase,
	}
}

func (we *WeatherByCepHttp) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	cep := r.PathValue("cep")
	input := usecase.InputWbcUsecase{Cep: cep}

	output := we.usecase.Execute(input)
	if output.Err != nil {
		err := output.Err
		if err == usecase.CanNotFindLocation {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(err.Error()))
			return
		}

		if err == entity.InvalidCep {
			w.WriteHeader(http.StatusUnprocessableEntity)
			w.Write([]byte(err.Error()))
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Println(output)
	json.NewEncoder(w).Encode(output)
}
