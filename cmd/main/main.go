package main

import (
	"context"
	"errors"
	"fmt"
	audioHl "github.com/TeaStealers-backend-sem4/internal/pkg/audio/delivery"
	audioUc "github.com/TeaStealers-backend-sem4/internal/pkg/audio/usecase"
	"github.com/TeaStealers-backend-sem4/internal/pkg/logger"
	"github.com/TeaStealers-backend-sem4/internal/pkg/middleware"
	wordH "github.com/TeaStealers-backend-sem4/internal/pkg/words/delivery"
	wordUc "github.com/TeaStealers-backend-sem4/internal/pkg/words/usecase"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	logr := logger.NewSlogLogger("log.txt")

	_ = godotenv.Load()
	r := mux.NewRouter().PathPrefix("/api").Subrouter()
	r.Use(middleware.CORSMiddleware, middleware.AccessLogMiddleware)
	r.HandleFunc("/ping", pingPongHandler).Methods(http.MethodGet)

	aUc := audioUc.NewAudioUsecase()
	auHandler := audioHl.NewAudioHandler(aUc, logr)

	wUc := wordUc.NewAudioUsecase()
	wHandler := wordH.NewWordHandler(wUc, logr)
	word := r.PathPrefix("/word").Subrouter()
	word.Handle("/get_word/{word}", http.HandlerFunc(wHandler.GetWord)).Methods(http.MethodGet)
	audio := r.PathPrefix("/audio").Subrouter()
	audio.Handle("/save_audio", http.HandlerFunc(auHandler.SaveAudio)).Methods(http.MethodPost, http.MethodOptions)
	audio.Handle("/translate_audio", http.HandlerFunc(auHandler.TranslateAudio)).Methods(http.MethodPost, http.MethodOptions)

	srv := &http.Server{
		Addr:              ":8080",
		Handler:           r,
		ReadHeaderTimeout: 10 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
	}

	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		fmt.Printf("Start server on %s\n", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			fmt.Printf("listen: %s\n", err)
		}
	}()

	sig := <-signalCh
	fmt.Printf("Received signal: %v\n", sig)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		fmt.Printf("Server shutdown failed: %s\n", err)
	}
}

func pingPongHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	//w.Write([]byte("Hello World"))
	fmt.Fprintln(w, "pong")
}
