/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log"
	"net/http"

	"github.com/bruno-holanda15/weather_by_cep_goexpert/configs"
	"github.com/bruno-holanda15/weather_by_cep_goexpert/internal/domain/entity"
	"github.com/bruno-holanda15/weather_by_cep_goexpert/internal/domain/usecase"
	"github.com/bruno-holanda15/weather_by_cep_goexpert/internal/infra/web"
	"github.com/spf13/cobra"
)

// wbcCmd represents the wbc command
var wbcCmd = &cobra.Command{
	Use:   "wbc",
	Short: "Return weather by cep",
	Long:  "Search weather by informing a specific cep.",
	Run:   StartWbc,
}

func StartWbc(cmd *cobra.Command, args []string) {
	fmt.Println("wbc called")

	loader := &configs.Loader{}
	loader.LoadEnv()

	infoSearcher := entity.NewInfosSearcher()
	weatherByCEPUsecase := usecase.NewWeatherByCepUsecase(infoSearcher)
	wbcHandler := web.NewWeatherByCepHttp(weatherByCEPUsecase)

	http.HandleFunc("/weather/{cep}", wbcHandler.FindTemps)

	fmt.Println("Listening http server http://localhost:8082")
	if err := http.ListenAndServe(":8082", nil); err != nil {
		log.Fatalf("error starting http server - %v", err)
	}
}

func init() {
	rootCmd.AddCommand(wbcCmd)
}
