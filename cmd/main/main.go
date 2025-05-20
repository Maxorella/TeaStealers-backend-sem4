package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	audioHl "github.com/TeaStealers-backend-sem4/internal/audio/delivery"
	moduleH "github.com/TeaStealers-backend-sem4/internal/module/delivery"
	moduleRep "github.com/TeaStealers-backend-sem4/internal/module/repo"
	moduleUc "github.com/TeaStealers-backend-sem4/internal/module/usecase"
	statHl "github.com/TeaStealers-backend-sem4/internal/stat/delivery"
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

	authH "github.com/TeaStealers-backend-sem4/internal/auth/delivery"
	authR "github.com/TeaStealers-backend-sem4/internal/auth/repo"
	authUc "github.com/TeaStealers-backend-sem4/internal/auth/usecase"
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
	r.Use(middleware2.RequestIDMiddleware, middleware2.CORSMiddleware, accessLogMiddleware)

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
	statRepo := statRep.NewRepository(db, logr)

	statUsecase := statUc.NewStatUsecase(statRepo, wRepo, logr)
	wordUsecase := wordUc.NewWordUsecase(wRepo, statRepo, logr)

	audioHandler := audioHl.NewAudioHandler(statUsecase, cfg, logr)
	statHandler := statHl.NewStatHandler(statUsecase, cfg, logr, minioStorageClient)
	wordHandler := wordH.NewWordHandler(wordUsecase, statUsecase, cfg, logr, minioStorageClient)

	modulRep := moduleRep.NewRepository(db, logr)
	modulUc := moduleUc.NewModuleUsecase(modulRep, logr)
	modulHandler := moduleH.NewModuleHandler(modulUc, cfg, logr)
	audio := r.PathPrefix("/audio").Subrouter()
	audio.Handle("/translate_audio", http.HandlerFunc(audioHandler.TranslateAudio)).Methods(http.MethodPost, http.MethodOptions)

	authRepo := authR.NewRepository(db)
	authUsecase := authUc.NewAuthUsecase(authRepo)
	autHandler := authH.NewAuthHandler(authUsecase)

	r.HandleFunc("/register", autHandler.SignUp).Methods(http.MethodPost, http.MethodOptions)
	r.HandleFunc("/login", autHandler.Login).Methods(http.MethodPost, http.MethodOptions)
	r.Handle("/logout", middleware2.JwtMiddleware(http.HandlerFunc(autHandler.Logout), authRepo)).Methods(http.MethodGet, http.MethodOptions)
	r.HandleFunc("/check_auth", autHandler.CheckAuth).Methods(http.MethodGet, http.MethodOptions)

	r.Handle("/current-word-module",
		middleware2.JwtMiddlewareOptional(http.HandlerFunc(wordHandler.GetCurrentModuleWordHandler), authRepo)).Methods(http.MethodGet)

	r.Handle("/current-phrase-module",
		middleware2.JwtMiddlewareOptional(http.HandlerFunc(wordHandler.GetCurrentModulePhraseHandler), authRepo)).Methods(http.MethodGet)

	r.Handle("/create-word-module", http.HandlerFunc(modulHandler.CreateModuleWordHandler)).Methods(http.MethodPost)
	r.Handle("/create-phrase-module", http.HandlerFunc(modulHandler.CreateModulePhraseHandler)).Methods(http.MethodPost)
	r.Handle("/word-exercises", http.HandlerFunc(wordHandler.CreateWordExerciseHandler)).Methods(http.MethodPost)
	r.Handle("/phrases-exercises", http.HandlerFunc(wordHandler.CreatePhraseExerciseHandler)).Methods(http.MethodPost)

	r.Handle("/exercise-progress",
		middleware2.JwtMiddleware(http.HandlerFunc(wordHandler.UpdateProgressHandler), authRepo)).Methods(http.MethodPost)

	r.Handle("/word-modules", http.HandlerFunc(wordHandler.WordModulesHandler)).Methods(http.MethodGet)
	r.Handle("/phrase-modules", http.HandlerFunc(wordHandler.PhraseModulesHandler)).Methods(http.MethodGet)

	r.Handle("/word-modules/{id}/exercises",
		middleware2.JwtMiddlewareOptional(http.HandlerFunc(wordHandler.GetWordModuleExercisesHandler), authRepo)).Methods(http.MethodGet)

	r.Handle("/phrase-modules/{id}/exercises",
		middleware2.JwtMiddlewareOptional(http.HandlerFunc(wordHandler.GetPhraseModuleExercisesHandler), authRepo)).Methods(http.MethodGet)

	word := r.PathPrefix("/word").Subrouter()
	topic := r.PathPrefix("/topic").Subrouter()
	tip := r.PathPrefix("/tip").Subrouter()
	//	user.Handle("/me", middleware.JwtMiddleware(http.HandlerFunc(userHandler.GetCurUser), authRepo)).Methods(http.MethodGet, http.MethodOptions)

	//word.Handle("/create_word", http.HandlerFunc(wordHandler.CreateWord)).Methods(http.MethodPost)
	word.Handle("/words_with_topic", http.HandlerFunc(wordHandler.WordsWithTopicHandler)).Methods(http.MethodPost)
	word.Handle("/{word}", http.HandlerFunc(wordHandler.GetWord)).Methods(http.MethodGet)
	topic.Handle("/all_topics", http.HandlerFunc(statHandler.GetAllTopics)).Methods(http.MethodGet)
	topic.Handle("/topic_progress", http.HandlerFunc(wordHandler.GetTopicProgressHandler)).Methods(http.MethodGet)
	tip.Handle("/get_tip", http.HandlerFunc(wordHandler.GetTipHandler)).Methods(http.MethodPost)
	tip.Handle("/upload_tip", http.HandlerFunc(wordHandler.UploadTipHandler)).Methods(http.MethodPost)
	// tip := r.PathPrefix("/tip").Subrouter()
	word.Handle("/rand/word", http.HandlerFunc(wordHandler.GetRandomWord)).Methods(http.MethodPost)
	/*
		word.Handle("/get_tags", http.HandlerFunc(wHandler.SelectTags)).Methods(http.MethodGet)
		word.Handle("/words_with_tag", http.HandlerFunc(wHandler.SelectWordsWithTopic)).Methods(http.MethodPost)

		tip.Handle("/upload_tip", http.HandlerFunc(wHandler.UploadTip)).Methods(http.MethodPost)
		tip.Handle("/get_tip", http.HandlerFunc(wHandler.GetTip)).Methods(http.MethodPost)
	*/
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
	fmt.Fprintln(w, "pong")
}
