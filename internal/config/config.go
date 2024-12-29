package config

import (
	"log"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type Config struct {
	Logger struct {
		Level      string `yaml:"level"`
		Filename   string `yaml:"filename"`
		MaxSize    int    `yaml:"max_size"`
		MaxBackups int    `yaml:"max_backups"`
		MaxAge     int    `yaml:"max_age"`
		Compress   bool   `yaml:"compress"`
	} `yaml:"logger"`

	Transport struct {
		HTTP struct {
			Host string `yaml:"host"`
			Port string `yaml:"port"`
		} `yaml:"http"`
	} `yaml:"transport"`

	External struct {
		Ethereum struct {
			RPCEndpoint string `env:"CRYPTOSERVICE_ETHEREUM_RPCENDPOINT"`
		} `yaml:"Ethereum"`
		Tron struct {
			RPCEndpoint string `env:"CRYPTOSERVICE_TRON_RPCENDPOINT"`
		} `yaml:"Tron"`
	} `yaml:"external"`

	Storages struct {
		Cache struct {
			Host             string `env:"CRYPTOSERVICE_CACHE_HOST"`
			Port             string `env:"CRYPTOSERVICE_CACHE_PORT"`
			Password         string `env:"CRYPTOSERVICE_CACHE_PASSWORD"`
			DBIndex          int    `yaml:"db_index"`
			WalletBalanceTTL int64  `yaml:"wallet_balance_ttl"`
		} `yaml:"cache"`
	} `yaml:"storages"`
}

func MustLoad() *Config {
	var cfg Config
	if _, err := os.Stat(".env"); err == nil {
		err := godotenv.Load()
		if err != nil {
			log.Fatal("error loading .env file")
		}
	}

	configPath := "./configs/app.yml"
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file does not exists: %s", configPath)
	}

	err := cleanenv.ReadConfig(configPath, &cfg)
	if err != nil {
		log.Fatalf("config file cannot be readed: %s", err)
	}

	return &cfg
}
