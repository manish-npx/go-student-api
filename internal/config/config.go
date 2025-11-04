package config

import (
	"flag"
	"log"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type HttpServer struct {
	Addr string `yaml:"address" env:"HTTP_ADDRESS" env-required:"true"`
}

type Postgres struct {
	Host     string `yaml:"host" env:"PG_HOST" env-required:"true"`
	Port     int    `yaml:"port" env:"PG_PORT" env-required:"true"`
	User     string `yaml:"user" env:"PG_USER" env-required:"true"`
	Password string `yaml:"password" env:"PG_PASSWORD" env-required:"true"`
	DBName   string `yaml:"dbname" env:"PG_DBNAME" env-required:"true"`
	SSLMode  string `yaml:"sslmode" env:"PG_SSLMODE" env-default:"disable"`
}

type Config struct {
	Env         string     `yaml:"env" env:"ENV" env-required:"true"`
	StoragePath string     `yaml:"storage_path" env:"STORAGE_PATH"`
	HttpServer  HttpServer `yaml:"http_server"`
	Postgres    Postgres   `yaml:"postgres"`
}

func MustLoad() *Config {
	var configPath string

	configPath = os.Getenv("CONFIG_PATH")
	if configPath == "" {
		flags := flag.String("config", "", "Path to configuration file")
		flag.Parse()

		configPath = *flags

		if configPath == "" {
			log.Fatal("Config Path not available")
		}
	}
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file does not exits: %s", configPath)
	}

	var cfg Config

	err := cleanenv.ReadConfig(configPath, &cfg)
	if err != nil {
		log.Fatalf("can not read config file : %s", err.Error())
	}

	return &cfg
}
