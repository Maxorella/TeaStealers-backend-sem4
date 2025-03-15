package config

import (
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	_ "github.com/joho/godotenv/autoload" // Load enviroment from .env
)

type Config struct {
	MlServer     MlService
	AudioUserDir string `env:"USER_AUDIO_DIR" env-default:"/ouzi/audio"`
}

type MlService struct {
	Address string        `env:"ML_ADDRESS" env-default:"178.57.232.224"`
	Port    string        `env:"ML_PORT" env-default:"5000"`
	Timeout time.Duration `env:"ML_TIMEOUT" env-default:"10s"`
}

func MustLoad() *Config {
	var cfg Config

	if err := cleanenv.ReadEnv(&cfg); err != nil {
		log.Printf("cannot read .env file: %s\n (fix: you need to put .env file in main dir)", err)
		os.Exit(1)
	}
	return &cfg
}
