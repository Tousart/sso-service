package config

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	// Env            string     `yaml:"env" env-default:"local"`
	// StoragePath    string     `yaml:"storage_path" env-required:"true"`
	GRPC           GRPCConfig `yaml:"grpc"`
	MigrationsPath string
	PSQL           Postgres `yaml:"psql"`
	// TokenTTL       time.Duration `yaml:"token_ttl" env-default:"1h"`
}

type GRPCConfig struct {
	Port    int           `yaml:"port" env:"GRPC_PORT"`
	Timeout time.Duration `yaml:"timeout" env:"GRPC_TIMEOUT"`
}

type Postgres struct {
	Host    string `yaml:"host" env:"PSQL_HOST"`
	Port    int    `yaml:"port" env:"PSQL_PORT"`
	DBName  string `yaml:"psql_name" env:"PSQL_NAME"`
	SSLMode string `yaml:"sslmode" env:"PSQL_SSLMODE"`
}

func parseFlags() string {
	cfgPathPtr := flag.String("config", "", "Path to cfg")
	flag.Parse()

	cfgPath := *cfgPathPtr

	return cfgPath
}

func MustLoad() (*Config, error) {
	cfgPath := parseFlags()

	if cfgPath == "" {
		return nil, errors.New("config path is empty")
	}

	if _, err := os.Stat(cfgPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("config does not exist: %v", err)
	}

	var cfg Config
	if err := cleanenv.ReadConfig(cfgPath, &cfg); err != nil {
		log.Fatal("error reading config")
		return nil, fmt.Errorf("failed to read config: %v", err)
	}

	if err := cleanenv.ReadEnv(cfg); err != nil {
		return nil, fmt.Errorf("failed to read env file: %v", err)
	}

	return &cfg, nil
}
