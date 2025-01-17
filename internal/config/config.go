package config

import (
	"errors"
	"flag"
	"strconv"
	"strings"
)

type Config struct {
	Host    string
	Port    int
	Address string
	BaseURL string
}

func (c *Config) String() string {
	return c.Host + ":" + strconv.Itoa(c.Port)
}

func (c *Config) Set(flagValue string) error {
	hp := strings.Split(flagValue, ":")
	if len(hp) != 2 {
		return errors.New("invalid flag value")
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
	config := &Config{
		Host:    "localhost",
		Port:    8080,
		Address: "localhost:8080",
	}
	flag.Var(config, "a", "host:port")
	flag.StringVar(&config.BaseURL, "b", "http://localhost:8080/", "base url")
	flag.Parse()

	return config
}
