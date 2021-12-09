package main

import (
	"github.com/mdbraber/acmeproxy/acmeproxy/server"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var (
	version = "dev"
)

func main() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("/etc/acmeproxy/")
	viper.AddConfigPath("$HOME/.acmeproxy")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()

	config := NewDefaultConfig()

	if _, ok := err.(viper.ConfigFileNotFoundError); ok {
		log.Info("no config file found")
	}

	err = viper.Unmarshal(&config)
	if err != nil {
		log.Fatalf("unable to decode config, %v", err)
	}

	err = setEnvVars(config.Environment)
	if err != nil {
		log.Fatalf("unable to set env vars, %v", err)
	}

	serverConfig := server.Config{
		HtpasswdFile:   config.Server.Htpasswd,
		AllowedIPs:     config.Filter.Ips,
		AllowedDomains: config.Filter.Domains,
		AccessLogFile:  config.Server.Accesslog,
		Port:           config.Server.Port,
		ProviderName:   config.Provider,
	}

	srv, err := server.NewServer(&serverConfig)
	if err != nil {
		log.Fatalf("unable to create server, %v", err)
	}

	srv.Run()
}
