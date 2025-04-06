package config

import (
	"flag"
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env string `yaml:"env" env-default:"prod"`
	//storage
	HTTPServer HTTPServer `yaml:"http_server"`
}

type HTTPServer struct {
	Addr        string        `yaml:"addr" env-required:"true"`
	Timeout     time.Duration `yaml:"timeout" end-default:"4s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"60s"`
}

func MustLoad() *Config {
	configPath := flag.String("config_path", "config_path", "path to config file")
	flag.Parse()
	if *configPath == "" {
		log.Fatal("CONFIG_PATH is not set")
	}
	log.Println(*configPath)

	//check if configFile exitst
	if _, err := os.Stat(*configPath); os.IsNotExist(err) {
		log.Fatalf("config file does not exits: %s", *configPath)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(*configPath, &cfg); err != nil {
		log.Fatalf("cannot read config: %s", err)
	}

	return &cfg

}
