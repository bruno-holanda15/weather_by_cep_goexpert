package web

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/bruno-holanda15/weather_by_cep_goexpert/internal/domain/entity"
	"github.com/bruno-holanda15/weather_by_cep_goexpert/internal/domain/usecase"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
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
	tracer trace.Tracer
}

func NewValidateCepHttp(usecase *usecase.ValidateCepUsecase, tracer trace.Tracer) *ValidateCepHttp {
	return &ValidateCepHttp{
		usecase: usecase,
		tracer: tracer,
	}
}

func (v *ValidateCepHttp) ValidateCep(w http.ResponseWriter, r *http.Request) {
	carrier := propagation.HeaderCarrier(r.Header)
	ctx := r.Context()
	ctx = otel.GetTextMapPropagator().Extract(ctx, carrier)

	ctx, span := v.tracer.Start(ctx, "validateCep Span")
	defer span.End()
	
	time.Sleep(300*time.Millisecond)

	var input usecase.InputWbcUsecase
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil || input.Cep == "" {
		http.Error(w, "Invalid zipcode", http.StatusUnprocessableEntity)
		return
	}

	output := v.usecase.Execute(ctx, input)
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
