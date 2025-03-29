package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/TeaStealers-backend-sem4/internal/pkg/config"
	"github.com/TeaStealers-backend-sem4/internal/pkg/logger"
	"github.com/TeaStealers-backend-sem4/internal/pkg/middleware"
	minioS "github.com/TeaStealers-backend-sem4/internal/pkg/minio"
	minioH "github.com/TeaStealers-backend-sem4/internal/pkg/minio/delivery"
	"github.com/joho/godotenv"

	"github.com/gorilla/mux"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	_ = godotenv.Load()
	cfg := config.MustLoad()

	//logr := logger.NewSlogLogger("log.txt")
	logr := logger.NewSlogStdOutLogger()
	logr.LogDebug("Started minio client logger")

	r := mux.NewRouter().PathPrefix("/api").Subrouter()
	r.Use(middleware.CORSMiddleware, middleware.AccessLogMiddleware)
	r.HandleFunc("/ping", pingPongHandler).Methods(http.MethodGet)

	minClient := minioS.NewMinioClient(cfg)
	err := minClient.InitMinio()
	if err != nil {
		logr.LogDebug(err.Error())
		os.Exit(-1)
	}

	minH := minioH.NewMinioHandler(minClient, cfg, logr)
	minH.RegisterRoutes(r)

	srv := &http.Server{
		Addr:              ":8081",
		Handler:           r,
		ReadHeaderTimeout: 10 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
	}

	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		str := fmt.Sprintf("Start server on %s\n", srv.Addr)
		logr.LogDebug(str)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			str := fmt.Sprintf("listen: %s\n", err)
			logr.LogDebug(str)
		}
	}()

	sig := <-signalCh
	str := fmt.Sprintf("Received signal: %v\n", sig)
	logr.LogDebug(str)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		str := fmt.Sprintf("Server shutdown failed: %s\n", err)
		logr.LogDebug(str)
	}
}

func pingPongHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	//w.Write([]byte("Hello World"))
	fmt.Fprintln(w, "pong")
}
