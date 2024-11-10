package config

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
)

type Config struct {
	Server struct {
		Level string `yaml:"level"`
		Host  string `yaml:"host"`
		Port  string `yaml:"port"`
	} `yaml:"server"`
	Database struct {
		PostgresSQL struct {
			Host              string `yaml:"host"`
			Port              string `yaml:"port"`
			Username          string `yaml:"username"`
			Password          string `yaml:"password"`
			Database          string `yaml:"database"`
			Url               string `yaml:"url"`
			MaxConnections    int32  `yaml:"max_connections"`
			MinConnections    int32  `yaml:"min_connections"`
			IdleTime          string `yaml:"idle_time"`
			HealthCheckPeriod string `yaml:"health_check_period"`
		} `yaml:"postgresql"`
		Redis struct {
		} `yaml:"redis"`
	} `yaml:"database"`
	JWT struct {
		JWTSecret         string `yaml:"jwt_secret"`
		JWTExpirationTime int    `yaml:"jwt_expiration_time"` // Время истечения токена в часах
	} `yaml:"jwt"`
}

func MustParse(configPath string) (*Config, error) {
	config := &Config{}
	file, err := os.Open(configPath)
	if err != nil {
		fmt.Printf("Error opening config file: %s\n", err)
		return nil, err
	}
	defer file.Close()

	d := yaml.NewDecoder(file)
	if err := d.Decode(config); err != nil {
		fmt.Printf("Error decoding config file: %s\n", err)
		return nil, err
	}

	return config, nil
}
