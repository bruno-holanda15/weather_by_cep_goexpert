/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"net/http"

	"github.com/bruno-holanda15/weather_by_cep_goexpert/internal/infra/web"
	"github.com/spf13/cobra"
)

// cepValidatorCmd represents the cepValidator command
var cepValidatorCmd = &cobra.Command{
	Use:   "cep_validator",
	Short: "Validate cep",
	Long:  "Validate cep to call the weather by cep service",
	Run:   StartValidator,
}

func StartValidator(cmd *cobra.Command, args []string) {
	fmt.Println("cepValidator called")

	http.HandleFunc("/weather", web.ValidateCep)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("toptop demais"))
	})

	http.ListenAndServe(":8081", nil)

}

func init() {
	rootCmd.AddCommand(cepValidatorCmd)
}
