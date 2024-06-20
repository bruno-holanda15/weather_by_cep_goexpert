/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/bruno-holanda15/weather_by_cep_goexpert/internal/domain/usecase"
	"github.com/bruno-holanda15/weather_by_cep_goexpert/internal/infra/web"
	o "github.com/bruno-holanda15/weather_by_cep_goexpert/pkg/otel"
	"github.com/spf13/cobra"
	"go.opentelemetry.io/otel"
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

	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, os.Interrupt, syscall.SIGINT)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGINT)
	defer stop()

	shutdown, err := o.InitProvider("cepValidator", "otel-collector:4317")
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := shutdown(ctx); err != nil {
			log.Fatal("failed to shutdown TracerProvider: %w", err)
		}
	}()

	tracer := otel.Tracer("cepValidator-tracer")

	validateCepUsecase := &usecase.ValidateCepUsecase{}
	validateCepHandler := web.NewValidateCepHttp(validateCepUsecase, tracer)

	http.HandleFunc("/weather", validateCepHandler.ValidateCep)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello cep validator!"))
	})

	server := &http.Server{
		Addr: ":8081",
	}

	go func() {
		fmt.Println("Listening http server http://localhost:8081")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Error starting http server port 8081: %v", err)
		}
	}()

	select {
	case <-signalCh:
		log.Println("Shutting down gracefully wbc1, CTRL+C pressed...")
	case <-ctx.Done():
		log.Println("Shutting down gracefully wbc1, interrupet system...")
	}

	shutdownContext, shutDownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutDownCancel()

	fmt.Println("Server shutting down wbc1")
	if err := server.Shutdown(shutdownContext); err != nil {
		log.Fatal("Error shutting down")
	}
}

func init() {
	rootCmd.AddCommand(cepValidatorCmd)
}
