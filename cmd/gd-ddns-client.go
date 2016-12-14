package main

import (
	"flag"
	"log"

	"github.com/twexler/gd-ddns-client"
)

var (
	configFile string
)

func init() {
	flag.StringVar(&configFile, "configFile", "/etc/gd-ddns-client.yaml", "spacify the config file to be used")
}

func main() {
	conf, err := gdDDNSClient.NewConfigFromFile(configFile)
	if err != nil {
		log.Fatalf("Unable to read config file %s: %s", configFile, err.Error())
	}
}
