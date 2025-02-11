package config

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

type AddressConfig struct {
	Host    string
	Port    int
	Address string
}

type BaseURLConfig struct {
	BaseURL string
}

type FileStorageConfig struct {
	Path string
}

type Config struct {
	A AddressConfig
	B BaseURLConfig
	F FileStorageConfig
}

func (c *AddressConfig) String() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}

func (c *AddressConfig) Set(flagValue string) error {
	hp := strings.Split(flagValue, ":")
	if len(hp) != 2 {
		return fmt.Errorf("invalid flag value")
	}
	port, err := strconv.Atoi(hp[1])
	if err != nil {
		return err
	}

	c.Host = hp[0]
	c.Port = port
	c.Address = c.String()
	return nil
}

func ParseFlags() *Config {
	addressCfg := &AddressConfig{
		Host:    "localhost",
		Port:    8080,
		Address: "localhost:8080",
	}

	baseURLCfg := &BaseURLConfig{}
	fileStorageCfg := &FileStorageConfig{}

	flag.Var(addressCfg, "a", "host:port (default: localhost:8080)")
	flag.StringVar(&baseURLCfg.BaseURL, "b", "http://localhost:8080", "base URL")
	flag.StringVar(&fileStorageCfg.Path, "f", "/tmp/short-url-db.json", "file storage path")
	flag.Parse()

	config := &Config{
		A: *addressCfg,
		B: *baseURLCfg,
		F: *fileStorageCfg,
	}

	if address, exists := os.LookupEnv("SERVER_ADDRESS"); exists {
		if err := config.A.Set(address); err != nil {
			log.Fatalf("Invalid SERVER_ADDRESS: %v", err)
		}
	}

	if baseURL, exists := os.LookupEnv("BASE_URL"); exists {
		config.B.BaseURL = baseURL
	}

	if fileStoragePath, exists := os.LookupEnv("FILE_STORAGE_PATH"); exists {
		config.F.Path = fileStoragePath
	}

	return config
}
