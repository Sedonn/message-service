package config

import (
	"flag"
	"os"
	"slices"

	"github.com/ilyakaznacheev/cleanenv"
)

// Все поддерживаемые типы окружения.
const (
	EnvLocal      = "local"
	EnvProduction = "production"
)

// Config хранит конфигурацию приложения.
type Config struct {
	Env  string     `yaml:"env" env-default:"local"`
	REST RESTConfig `yaml:"rest"`
	DB   DBConfig   `yaml:"db"`
}

// RESTConfig хранит конфигурацию REST-API сервера.
type RESTConfig struct {
	Port int `yaml:"port" env:"REST_PORT"`
}

// DBConfig хранит конфигурацию подключения к базе данных.
type DBConfig struct {
	Host     string `yaml:"host" env:"DB_HOST" env-required:"true"`
	Port     int    `yaml:"port" env:"DB_PORT" env-required:"true"`
	User     string `yaml:"username" env:"DB_USER" env-required:"true"`
	Password string `yaml:"password" env:"DB_PASSWORD" env-required:"true"`
	Database string `yaml:"database" env:"DB_NAME" env-required:"true"`
}

// MustLoad загружает текущую конфигурацию микросервиса на основе пути к файлу конфигурации,
// получаемого из флага запуска или переменной окружения.
//
// Конфигурация грузится по приоритету:
// 1. Флаг.
// 2. Переменная окружения.
func MustLoad() *Config {
	return MustLoadByPath(getConfigPath())
}

// MustLoad загружает текущую конфигурацию микросервиса на основе переданного пути к файлу конфигурации.
//
// Конфигурация грузится по приоритету:
// 1. Флаг.
// 2. Переменная окружения.
func MustLoadByPath(path string) *Config {
	if path == "" {
		panic("config path not set")
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		panic("config file not found by path: " + path)
	}

	var cfg Config
	if err := cleanenv.ReadConfig(path, &cfg); err != nil {
		panic("failed to load config: " + err.Error())
	}

	if !validateEnv(cfg.Env) {
		panic("unknown env: " + cfg.Env)
	}

	return &cfg
}

// getConfigPath получает путь к файлу конфигурации со следующим приоритетом:
// 1. Флаг.
// 2. Переменная окружения.
func getConfigPath() string {
	const configPathKey = "CONFIG_PATH"

	var path string
	flag.StringVar(&path, "config_path", "", "path to the .yaml config file.")
	flag.Parse()

	if path == "" {
		path = os.Getenv(configPathKey)
	}

	return path
}

// validateEnv выполняет валидацию текущей конфигурации.
func validateEnv(env string) bool {
	envTypes := []string{EnvLocal, EnvProduction}

	return slices.Contains(envTypes, env)
}
