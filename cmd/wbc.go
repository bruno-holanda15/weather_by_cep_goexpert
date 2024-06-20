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

	"github.com/bruno-holanda15/weather_by_cep_goexpert/configs"
	"github.com/bruno-holanda15/weather_by_cep_goexpert/internal/domain/entity"
	"github.com/bruno-holanda15/weather_by_cep_goexpert/internal/domain/usecase"
	"github.com/bruno-holanda15/weather_by_cep_goexpert/internal/infra/web"
	o "github.com/bruno-holanda15/weather_by_cep_goexpert/pkg/otel"
	"github.com/spf13/cobra"
	"go.opentelemetry.io/otel"
)

// wbcCmd represents the wbc command
var wbcCmd = &cobra.Command{
	Use:   "wbc",
	Short: "Return weather by cep",
	Long:  "Search weather by informing a specific cep, we created a shortcut/acronym wbc :)",
	Run:   StartWbc,
}

func StartWbc(cmd *cobra.Command, args []string) {
	fmt.Println("wbc called")

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

	loader := &configs.Loader{}
	loader.LoadEnv()

	infoSearcher := entity.NewInfosSearcher()
	weatherByCEPUsecase := usecase.NewWeatherByCepUsecase(infoSearcher)
	wbcHandler := web.NewWeatherByCepHttp(weatherByCEPUsecase, tracer)

	http.HandleFunc("/weather/{cep}", wbcHandler.FindTemps)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello weather by cep!"))
	})

	server := &http.Server{
		Addr: ":8082",
	}

	go func() {
		fmt.Println("Listening http server http://localhost:8082")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Error starting http server port 8082: %v", err)
		}
	}()

	select {
	case <-signalCh:
		log.Println("Shutting down gracefully wbc2, CTRL+C pressed...")
	case <-ctx.Done():
		log.Println("Shutting down gracefully wbc2, interrupet system...")
	}

	shutdownContext, shutDownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutDownCancel()

	fmt.Println("Server shutting down wbc2")
	if err := server.Shutdown(shutdownContext); err != nil {
		log.Fatal("Error shutting down")
	}
}

func init() {
	rootCmd.AddCommand(wbcCmd)
}
