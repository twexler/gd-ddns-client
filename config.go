package gdDDNSClient

import (
	"fmt"
	"io/ioutil"
	"time"

	yaml "gopkg.in/yaml.v2"
)

const (
	defaultDomainsURL     = "https://domains.google.com"
	defaultUpdateInterval = time.Hour
)

// Config is an interface that all config types must implement
type Config interface {
	// GetAllDomains returns a list of domains stored in the configuration struct
	GetAllDomains() []string
	// GetCredentialsByDomain returns an individual Credential for the specified domains
	GetCredentialsByHostname(hostname string) (Credential, error)
	// GetDomainsURL should return the URL for Google Domains' api
	GetDomainsURL() string
	// GetIPIfyURL returns the IPify URL to use in the client
	GetIPIfyURL() string
	// GetUpdateInterval returns the interval for the client to update Google Domains
	GetUpdateInterval() time.Duration
}

type config struct {
	Credentials    map[string]Credential
	DomainsURL     string        `yaml:"domainsURL"`
	IPIfyURL       string        `yaml:"ipifyURL"`
	UpdateInterval time.Duration `yaml:"updateInterval"`
}

// Credential is a simple struct to handle username/password combos
type Credential struct {
	User, Password string
}

// NewConfigFromFile creates a new Config from filename or returns any errors during file reads or YAML parsing
func NewConfigFromFile(filename string) (Config, error) {
	var conf config
	var confBytes []byte
	var err error
	if confBytes, err = ioutil.ReadFile(filename); err != nil {
		return &conf, err
	}
	if err = yaml.Unmarshal(confBytes, &conf); err != nil {
		return &conf, err
	}
	return &conf, nil
}

func (c config) GetAllDomains() []string {
	ret := make([]string, 0, len(c.Credentials))
	for k := range c.Credentials {
		ret = append(ret, k)
	}
	return ret
}

func (c config) GetCredentialsByHostname(hostname string) (Credential, error) {
	if v, ok := c.Credentials[hostname]; ok {
		return v, nil
	}
	return Credential{}, fmt.Errorf("No credentials for %s", hostname)
}

func (c config) GetDomainsURL() string {
	if c.DomainsURL == "" {
		return defaultDomainsURL
	}
	return c.DomainsURL
}

func (c config) GetIPIfyURL() string {
	if c.IPIfyURL == "" {
		return defaultIPIfyURL
	}
	return c.IPIfyURL
}

func (c config) GetUpdateInterval() time.Duration {
	if c.UpdateInterval == 0 {
		return defaultUpdateInterval
	}
	return c.UpdateInterval
}
