package main

import (
	"flag"
	"net"
	"time"

	"github.com/uber-go/zap"

	client "github.com/twexler/gd-ddns-client"
)

var (
	configFile string
	logger     zap.Logger
	currentIP  net.IP
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
	logger.Info("Using config file", zap.String("configFile", configFile))
	ticker := time.NewTicker(conf.GetUpdateInterval())
	ipifyClient := client.NewIPIfyAPI(conf.GetIPIfyURL())
	for _, hostname := range conf.GetAllDomains() {
		if cred, err := conf.GetCredentialsByHostname(hostname); err == nil {
			domainsClient := client.NewDomainsAPI(conf.GetDomainsURL())
			go runUpdates(ticker.C, hostname, cred, domainsClient, ipifyClient)
		} else {
			// this should be impossible, but who knows.
			logger.Error("Unable to get credentials for hostname", zap.String("hostname", hostname), zap.Error(err))
		}
	}
	// block!
	select {}
}

func runUpdates(
	c <-chan time.Time,
	hostname string,
	creds client.Credential,
	domainsClient client.DomainsAPI,
	ipifyClient client.IPIfyAPI) {
	for range c {
		if addr, err := ipifyClient.GetIPAddress(); err == nil {
			if !addr.Equal(currentIP) {
				logger.Info("Public IP address changed, updating", zap.String("old-ip", currentIP.String()), zap.String("new-ip", addr.String()))
				// eventually we might do something with offline, but I really don't understand what it's for right now
				if err := domainsClient.Update(creds, hostname, addr, false); err == nil {
					logger.Info("Successfully updated IP address for hostname", zap.String("hostname", hostname), zap.String("ip", addr.String()))
					currentIP = addr
					continue
				} else {
					logger.Error("Unable to update Google Domains", zap.Error(err))
				}
			} else {
				logger.Info("Public IP address not changed, not sending updates", zap.String("ip", currentIP.String()))
			}
		} else {
			logger.Error("Unable to get IP address from ipify", zap.Error(err))
		}
	}
}
