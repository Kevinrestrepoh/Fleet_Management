package main

import (
	"log"

	"github.com/Kevinrestrepoh/control-api/grpc"
	"github.com/Kevinrestrepoh/control-api/http"
	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()

	ingestion, err := grpc.NewIngestionClient("localhost:50051")
	if err != nil {
		log.Fatal(err)
	}

	handler := http.NewHandler(ingestion)

	e.POST("/vehicles/:id/command", handler.SendCommand)
	e.GET("/metrics/stream", handler.MetricsSSE)

	e.Logger.Fatal(e.Start(":8080"))
}
