package main

import (
	log "github.com/sirupsen/logrus"
	"github.com/walnuss0815/acmeproxy/v2/acme/provider"
	"github.com/walnuss0815/acmeproxy/v2/acme/proxy"
	"github.com/walnuss0815/acmeproxy/v2/server"

	"github.com/spf13/viper"
)

func main() {
	Init()

	port := viper.GetUint("server.port")
	allowedDomains := viper.GetStringSlice("allowed_domains")
	proxy := proxy.NewProxy(provider.NewDefaultProviderCloudflare(), allowedDomains)

	srv := server.NewServer(port, proxy)

	srv.Run()
}

func Init() {
	viper.SetEnvPrefix("acmeproxy")
	viper.AutomaticEnv()

	viper.SetConfigName("config")          // name of config file (without extension)
	viper.SetConfigType("yaml")            // REQUIRED if the config file does not have the extension in the name
	viper.AddConfigPath("/etc/acmeproxy/") // path to look for the config file in
	viper.AddConfigPath(".")               // optionally look for config in the working directory

	viper.SetDefault("server.port", "8080")

	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		log.Errorf("Could not read config file: %s", err.Error())
	} else {
		log.Info("Read config from file")
	}
}
