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

type Config struct {
	A AddressConfig
	B BaseURLConfig
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

	baseURlCfg := &BaseURLConfig{}

	flag.Var(addressCfg, "a", "host:port (default: localhost:8080)")
	flag.StringVar(&baseURlCfg.BaseURL, "b", "http://localhost:8080", "base URL")
	flag.Parse()

	config := &Config{
		A: *addressCfg,
		B: *baseURlCfg,
	}

	if address, exists := os.LookupEnv("SERVER_ADDRESS"); exists {
		if err := config.A.Set(address); err != nil {
			log.Fatalf("Invalid SERVER_ADDRESS: %v", err)
		}
	}

	if baseURL, exists := os.LookupEnv("BASE_URL"); exists {
		config.B.BaseURL = baseURL
	}

	return config
}
