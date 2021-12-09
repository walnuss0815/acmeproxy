package main

import (
	"github.com/mcuadros/go-defaults"
	"os"
	"strings"
)

type Config struct {
	Environment map[string]string
	Server      struct {
		Port      int    `default:"9096"`
		Htpasswd  string `default:""`
		Accesslog string `default:"/var/log/acmeproxy.log"`
	}
	Provider string
	Filter   struct {
		Ips     []string
		Domains []string
	}
}

func NewDefaultConfig() *Config {
	c := new(Config)
	defaults.SetDefaults(c)
	return c
}

func setEnvVars(vars map[string]string) error {
	for k, v := range vars {
		k = strings.ToUpper(k)
		err := os.Setenv(k, v)
		if err != nil {
			return err
		}
	}
	return nil
}
