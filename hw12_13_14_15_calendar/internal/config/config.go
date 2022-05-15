package config

import (
	"os"

	env "github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
)

type Config struct {
	Logger  LoggerConf
	Storage StorageConf
	Server  ServerConf
	Queue   QueueConf
}

type LoggerConf struct {
	Level string `env:"LOGGER_LEVEL" envDefault:"info"`
	File  string `env:"LOGGER_FILE" envDefault:"calendar.log"`
}

type StorageConf struct {
	Type string `env:"STORAGE_TYPE" envDefault:"notmemory"`
	Dsn  string `env:"STORAGE_DSN" envDefault:"postgres://postgres:postgres@postgres:5432/postgres"`
}

type ServerConf struct {
	Hport string `env:"HTTP_PORT" envDefault:"8080"`
	Gport string `env:"GRPC_PORT" envDefault:"50051"`
}

type QueueConf struct {
	URI string `env:"RABBIT_URI" envDefault:"amqp://guest:guest@rabbit:5672/"`
}

func LoadConfig(cfg interface{}, fileNames ...string) {
	if len(fileNames) == 0 {
		fileNames = []string{".env", ".env.local"}
	}

	valid := []string{""}
	for _, f := range fileNames {
		if _, err := os.Stat(f); os.IsNotExist(err) {
			continue
		}
		valid = append(valid, f)
	}
	_ = godotenv.Overload(valid...)
	_ = env.Parse(cfg)
}
