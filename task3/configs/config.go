package config

import (
	"encoding/json"
	"errors"
	"net"
	"os"
)

// simple struct for loading JSON config
type Config struct {
	Appname  string
	IP       string
	Port     string
	Username string
	Password string
}

// return config in a format for suitable ami
func (cfg *Config) GetConfig() (string, string, string, string) {
	return cfg.Appname, cfg.IP + ":" + cfg.Port, cfg.Username, cfg.Password
}

// MakeConfig sets up ari configuration from provided parameters
func MakeConfig(ip, port, user, pass string) (*Config, error) {
	if user == "" || pass == "" {
		return nil, errors.New("login credentials must be provided")
	}

	// if ip was not provided assume localhost
	if ip == "" {
		ip = "127.0.0.1"
	} else if net.ParseIP(ip) == nil { // if ip parser doesnt return a pointer to net.IP then an invalid adress string has been provided
		return nil, errors.New("invalid ip address")
	}

	if port == "" {
		port = "5038" // default asterisk port for telnet connection
	}

	return &Config{
		IP:       ip,
		Port:     port,
		Username: user,
		Password: pass,
	}, nil
}

// LoadConfig takes a json file as a parameter and loads configuration from it
func LoadConfig(filename string) (*Config, error) {
	config := &Config{}

	fileContent, err := os.ReadFile(filename)
	if err != nil {
		return config, err
	}

	err = json.Unmarshal(fileContent, &config)
	if err != nil {
		return config, err
	}

	return config, nil
}
