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
	GRPC      GRPCConfig     `yaml:"grpc"`
	Postgres  PostgresConfig `yaml:"postgres"`
	Redis     RedisConfig    `yaml:"redis"`
	Kafka     KafkaConfig    `yaml:"kafka"`
	JWTSecret string         `yaml:"jwt_secret" env:"JWT_SECRET"`

	// MigrationsPath string
	// Env            string     `yaml:"env" env-default:"local"`
	// StoragePath    string     `yaml:"storage_path" env-required:"true"`
	// TokenTTL       time.Duration `yaml:"token_ttl" env-default:"1h"`
}

type GRPCConfig struct {
	Port    int           `yaml:"port" env:"GRPC_PORT"`
	Timeout time.Duration `yaml:"timeout" env:"GRPC_TIMEOUT"`
}

type PostgresConfig struct {
	Host     string `yaml:"host" env:"POSTGRES_HOST"`
	Port     int    `yaml:"port" env:"POSTGRES_PORT"`
	DBName   string `yaml:"psql_db" env:"POSTGRES_DB"`
	User     string `yaml:"user" env:"POSTGRES_USER"`
	Password string `yaml:"password" env:"POSTGRES_PASSWORD"`
}

type RedisConfig struct {
	Host     string `yaml:"host" env:"REDIS_HOST"`
	Port     int    `yaml:"port" env:"REDIS_PORT"`
	Password string `yaml:"password" env:"REDIS_PASSWORD"`
	DB_ID    int    `yaml:"db_id" env:"DB_ID"`
}

type KafkaConfig struct {
	Brokers   string `yaml:"brokers" env:"KAFKA_BROKERS"`
	TopicName string `yaml:"topic_name" env:"KAFKA_TOPIC"`
	GroupID   string `yaml:"group_id" env:"KAFKA_GROUP"`
}

func ParseFlags() string {
	cfgPathPtr := flag.String("config", "", "Path to cfg")
	flag.Parse()

	cfgPath := *cfgPathPtr

	return cfgPath
}

func MustLoad(cfgPath string) (*Config, error) {
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

	if err := cleanenv.ReadEnv(&cfg); err != nil {
		return nil, fmt.Errorf("failed to read env file: %v", err)
	}

	return &cfg, nil
}
