package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	audioHl "github.com/TeaStealers-backend-sem4/internal/audio/delivery"
	statRep "github.com/TeaStealers-backend-sem4/internal/stat/repo"
	statUc "github.com/TeaStealers-backend-sem4/internal/stat/usecase"
	wordH "github.com/TeaStealers-backend-sem4/internal/word/delivery"
	wordRep "github.com/TeaStealers-backend-sem4/internal/word/repo"
	wordUc "github.com/TeaStealers-backend-sem4/internal/word/usecase"
	"github.com/TeaStealers-backend-sem4/pkg/config"
	"github.com/TeaStealers-backend-sem4/pkg/logger"
	middleware2 "github.com/TeaStealers-backend-sem4/pkg/middleware"
	minioS "github.com/TeaStealers-backend-sem4/pkg/minio"
	minioH "github.com/TeaStealers-backend-sem4/pkg/minio/delivery"
	utils2 "github.com/TeaStealers-backend-sem4/pkg/utils"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	_ = godotenv.Load()
	cfg := config.MustLoad()
	logr := logger.NewSlogStdOutLogger()
	logr.LogDebug("started slog")
	//logr := logger.NewSlogLogger("log.txt") TODO если хотим записывать в файл

	db, err := sql.Open("postgres", fmt.Sprintf("postgres://%v:%v@%v:%v/%v?sslmode=disable",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASS"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME")))
	if err != nil {
		panic("failed to connect database" + err.Error())
	}

	if err = db.Ping(); err != nil {
		log.Println("fail ping postgres")
		err = fmt.Errorf("error happened in db.Ping: %w", err)
		log.Println(err)
	}

	r := mux.NewRouter().PathPrefix("/api").Subrouter()
	accessLogMiddleware := middleware2.NewAccessLogMiddleware(logr)
	r.Use(middleware2.CORSMiddleware, middleware2.RequestIDMiddleware, accessLogMiddleware)

	r.HandleFunc("/ping", pingPongHandler).Methods(http.MethodGet)

	minClient := minioS.NewMinioClient(cfg, logr)
	err = minClient.InitMinio()
	if err != nil {
		logr.LogDebug(err.Error())
		os.Exit(-1)
	}

	minH := minioH.NewMinioHandler(minClient, cfg, logr)
	minH.RegisterRoutes(r)

	minioStorageClient := utils2.NewFileStorageClient(cfg.MinCli.AddressPort)

	wRepo := wordRep.NewRepository(db, logr)
	wUc := wordUc.NewWordUsecase(wRepo, logr)
	statRepo := statRep.NewRepository(db, logr)
	statUsecase := statUc.NewStatUsecase(statRepo, wRepo, logr)

	auHandler := audioHl.NewAudioHandler(statUsecase, cfg, logr)
	audio := r.PathPrefix("/audio").Subrouter()
	audio.Handle("/translate_audio", http.HandlerFunc(auHandler.TranslateAudio)).Methods(http.MethodPost, http.MethodOptions)

	wHandler := wordH.NewWordHandler(wUc, cfg, logr, minioStorageClient)
	word := r.PathPrefix("/word").Subrouter()
	tip := r.PathPrefix("/tip").Subrouter()
	word.Handle("/rand/word", http.HandlerFunc(wHandler.GetRandomWord)).Methods(http.MethodPost)
	word.Handle("/get_tags", http.HandlerFunc(wHandler.SelectTags)).Methods(http.MethodGet)
	word.Handle("/stat/write_stat", http.HandlerFunc(wHandler.WriteStat)).Methods(http.MethodPost)
	word.Handle("/stat/get_stat/{word_id}", http.HandlerFunc(wHandler.GetStat)).Methods(http.MethodGet)
	word.Handle("/words_with_tag", http.HandlerFunc(wHandler.SelectWordsWithTopic)).Methods(http.MethodPost)
	word.Handle("/{word}", http.HandlerFunc(wHandler.GetWord)).Methods(http.MethodGet)
	word.Handle("/create_word", http.HandlerFunc(wHandler.CreateWordHandler)).Methods(http.MethodPost)
	word.Handle("/pronunciation/{word}", http.HandlerFunc(wHandler.UploadAudioHandler)).Methods(http.MethodPost)
	tip.Handle("/upload_tip", http.HandlerFunc(wHandler.UploadTip)).Methods(http.MethodPost)
	tip.Handle("/get_tip", http.HandlerFunc(wHandler.GetTip)).Methods(http.MethodPost)
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
