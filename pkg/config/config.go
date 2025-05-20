package config

import (
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	_ "github.com/joho/godotenv/autoload" // Load enviroment from .env
)

type Config struct {
	MlServer        MlService
	AudioUserDir    string `env:"USER_AUDIO_DIR" env-default:"/ouzi/audio"`
	AudioExampleDir string `env:"EXAMPLE_AUDIO_DIR" env-default:"/ouzi/examples/"`
	MinioService    MinioS3
	MinCli          MinioClient
}

/*
type Database struct {
	DBName string `env:"POSTGRES_DB" env-required:"true"`
	DBPass string `env:"POSTGRES_PASSWORD" env-required:"true"`
	DBHost string `env:"DB_HOST" env-default:"0.0.0.0"`
	DBPort int    `env:"DB_PORT" env-required:"true"`
	DBUser string `env:"POSTGRES_USER" env-required:"true"`
}

*/

type MinioS3 struct {
	MinioEndpoint     string `env:"MINIO_ENDPOINT" env-default:"host.docker.internal:9000"`
	BucketName        string `env:"MINIO_BUCKET_NAME" env-default:"defaultbucket"`
	MinioRootUser     string `env:"MINIO_ROOT_USER" env-default:"minioadmin"`
	MinioRootPassword string `env:"MINIO_ROOT_PASSWORD" env-default:"minioadmin"`
	MinioUseSSL       bool   `env:"MINIO_USE_SSL" env-default:"false"`
}

type MinioClient struct {
	AddressPort string `env:"MINIO_CLIENT_ADR_PORT" env-default:"http://localhost:8080"`
}

type MlService struct {
	Address                  string        `env:"ML_ADDRESS" env-default:"178.57.232.224"`
	Port                     string        `env:"ML_PORT" env-default:"5000"`
	Timeout                  time.Duration `env:"ML_TIMEOUT" env-default:"10s"`
	TranscribeWordEndpoint   string        `env:"ML_ENDPOINT_WORD"`
	TranscribePhraseEndpoint string        `env:"ML_ENDPOINT_PHRASE"`
	HelpTextEndpoint         string        `env:"ML_ENDPOINT_HELPTEXT"`
	DialogEndpoint           string        `env:"ML_ENDPOINT_DIALOG"`
}

func MustLoad() *Config {
	var cfg Config

	if err := cleanenv.ReadEnv(&cfg); err != nil {
		log.Printf("cannot read .env file: %s\n (fix: you need to put .env file in main dir)", err)
		os.Exit(1)
	}
	return &cfg
}
