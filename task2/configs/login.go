package config

import (
	"errors"
	"net"
)

type Conf struct {
	IP       string
	Port     string
	Username string
	Secret   string
}

func NewConf(ip, port, user, secret string) (*Conf, error) {
	if user == "" || secret == "" {
		return nil, errors.New("login credentials must be provided")
	}

	if ip == "" {
		ip = "127.0.0.1"
	} else if net.ParseIP(ip) == nil {
		return nil, errors.New("invalid ip address")
	}

	if port == "" {
		port = "5038"
	}

	return &Conf{
		IP:       ip,
		Port:     port,
		Username: user,
		Secret:   secret,
	}, nil
}
