package main

import (
	"flag"

	"github.com/uber-go/zap"

	client "github.com/twexler/gd-ddns-client"
)

var (
	configFile string
	logger     zap.Logger
)

func init() {
	flag.StringVar(&configFile, "configFile", "/etc/gd-ddns-client.yaml", "spacify the config file to be used")
	logger = zap.New(zap.NewTextEncoder())
}

func main() {
	conf, err := client.NewConfigFromFile(configFile)
	if err != nil {
		logger.Fatal("Unable to read config file", zap.String("configFile", configFile), zap.Error(err))
	}

}
