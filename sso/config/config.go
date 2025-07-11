package config

import (
	"flag"
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
	// TokenTTL       time.Duration `yaml:"token_ttl" env-default:"1h"`
}

type GRPCConfig struct {
	Port    int           `yaml:"port"`
	Timeout time.Duration `yaml:"timeout"`
}

// Получение пути к конфиг-файлу
func parseFlags() string {
	cfgPathPtr := flag.String("config", "", "Path to cfg")
	flag.Parse()

	cfgPath := *cfgPathPtr

	// fmt.Println(cfgPath)

	// if cfgPath == "" {
	// 	cfgPath = os.Getenv("CONFIG_PATH")
	// }

	return cfgPath
}

// Парсинг конфиг-файла
func MustLoad() *Config {
	cfgPath := parseFlags()

	if cfgPath == "" {
		log.Fatal("config path is empty")
	}

	// Существует ли конфиг по пути
	if _, err := os.Stat(cfgPath); os.IsNotExist(err) {
		log.Fatal("config does not exist: " + cfgPath)
	}

	// Читаем файл конфига и заполняем cfg
	var cfg Config
	if err := cleanenv.ReadConfig(cfgPath, &cfg); err != nil {
		log.Fatal("error reading config")
	}

	return &cfg
}
