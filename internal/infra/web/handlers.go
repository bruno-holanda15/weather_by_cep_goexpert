package web

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/bruno-holanda15/weather_by_cep_goexpert/internal/domain/entity"
	"github.com/bruno-holanda15/weather_by_cep_goexpert/internal/domain/usecase"
)

type WeatherByCepHttp struct {
	usecase *usecase.WeatherByCepUsecase
}

func NewWeatherByCepHttp(usecase *usecase.WeatherByCepUsecase) *WeatherByCepHttp {
	return &WeatherByCepHttp{
		usecase: usecase,
	}
}

func (we *WeatherByCepHttp) FindTemps(w http.ResponseWriter, r *http.Request) {
	cep := r.PathValue("cep")
	input := usecase.InputWbcUsecase{Cep: cep}

	output := we.usecase.Execute(input)
	if output.Err != nil {
		err := output.Err
		if err == entity.ErrorCanNotFindLocation {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(err.Error()))
			return
		}

		if err == entity.ErrorInvalidCep {
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

type CepBody struct {
	Cep string `json:"cep"`
}

type ValidateCepHttp struct {
	usecase *usecase.ValidateCepUsecase
}

func NewValidateCepHttp(usecase *usecase.ValidateCepUsecase) *ValidateCepHttp {
	return &ValidateCepHttp{
		usecase: usecase,
	}
}

func (v *ValidateCepHttp) ValidateCep(w http.ResponseWriter, r *http.Request) {
	var input usecase.InputWbcUsecase
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil || input.Cep == "" {
		http.Error(w, "Invalid zipcode", http.StatusUnprocessableEntity)
		return
	}

	output := v.usecase.Execute(input)
	if err := output.Err; err != nil {
		if err == entity.ErrorInvalidCep {
			w.WriteHeader(http.StatusUnprocessableEntity)
			w.Write([]byte("invalid zipcode"))
			return
		}

		if err == entity.ErrorEmptyCep {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("can not find zipcode"))
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
