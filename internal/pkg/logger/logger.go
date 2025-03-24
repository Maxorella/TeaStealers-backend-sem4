package logger

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
)

const (
	DeliveryLayer   = "deliveryLayer"
	UsecaseLayer    = "usecaseLayer"
	RepositoryLayer = "repositoryLayer"
)

type Logger interface {
	LogDebug(message string)
	LogInfo(requestId string, layer string, methodName string, message string)
	LogError(requestId string, layer string, methodName string, err error)
	LogErrorResponse(requestId string, layer string, methodName string, err error, status int)
	LogSuccess(requestId string, layer string, methodName string)
	LogSuccessResponse(requestId string, layer string, methodName string)
}

type SlogLogger struct {
	logger *slog.Logger
}

func NewSlogStdOutLogger() *SlogLogger {
	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})

	return &SlogLogger{
		logger: slog.New(handler),
	}
}

func NewSlogLogger(logFile string) *SlogLogger {
	file, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		slog.Error("Failed to open log file", slog.String("error", err.Error()))
		os.Exit(1)
	}

	handler := slog.NewJSONHandler(file, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})

	return &SlogLogger{
		logger: slog.New(handler),
	}
}

func (l *SlogLogger) LogDebug(message string) {
	l.logger.Debug(
		fmt.Sprintf("MESSAGE: %v", message),
	)
}

func (l *SlogLogger) LogInfo(requestId string, layer string, methodName string, message string) {
	l.logger.Info(
		fmt.Sprintf("INFO: %v", message),
		slog.String("layer", layer),
		slog.String("method", methodName),
		slog.String("requestId", requestId),
	)
}

func (l *SlogLogger) LogError(requestId string, layer string, methodName string, err error) {
	l.logger.Error(
		fmt.Sprintf("ERROR: %v", err.Error()),
		slog.String("layer", layer),
		slog.String("method", methodName),
		slog.String("requestId", requestId),
	)
}

func (l *SlogLogger) LogErrorResponse(requestId string, layer string, methodName string, err error, status int) {
	l.logger.Error(
		fmt.Sprintf("ERROR RESPONSE: %v", err.Error()),
		slog.String("layer", layer),
		slog.String("method", methodName),
		slog.String("requestId", requestId),
		slog.Int("responseStatus", status),
	)
}

func (l *SlogLogger) LogSuccess(requestId string, layer string, methodName string) {
	l.logger.Info(
		fmt.Sprintf("OK: %v", requestId),
		slog.String("layer", layer),
		slog.String("method", methodName),
		slog.String("requestId", requestId),
	)
}

func (l *SlogLogger) LogSuccessResponse(requestId string, layer string, methodName string) {
	l.logger.Info(
		fmt.Sprintf("OK RESPONSE: %v", requestId),
		slog.String("layer", layer),
		slog.String("method", methodName),
		slog.String("requestId", requestId),
		slog.Int("responseStatus", http.StatusOK),
	)
}
