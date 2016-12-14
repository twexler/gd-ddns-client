package gdDDNSClient

import (
	"fmt"
	"io/ioutil"

	yaml "gopkg.in/yaml.v2"
)

var defaultDomainsURL = "https://domains.google.com"

// Config an interface that all config types must implement
type Config interface {
	// GetDomainsURL should return the URL for Google Domains' api
	GetDomainsURL() string
	// GetAllDomains returns a list of domains stored in the configuration struct
	GetAllDomains() []string
	// GetCredentialsByDomain returns an individual Credential for the specified domains
	GetCredentialsByDomain(domain string) (Credential, error)
}

type config struct {
	domainsURL  string
	credentials map[string]Credential
}

// Credential is a simple struct to handle username/password combos
type Credential struct {
	User, Password string
}

// NewConfig creates a new Config struct based on values provided
func NewConfig(domainsURL string, credentials map[string]Credential) Config {
	return config{
		domainsURL:  domainsURL,
		credentials: credentials,
	}
}

// NewConfigFromFile creates a new Config from filename or returns any errors during file reads or YAML parsing
func NewConfigFromFile(filename string) (Config, error) {
	var conf *config
	var confBytes []byte
	var err error
	if confBytes, err = ioutil.ReadFile(filename); err != nil {
		return conf, err
	}
	if err = yaml.Unmarshal(confBytes, conf); err != nil {
		return conf, err
	}
	return conf, nil
}

func (c config) GetDomainsURL() string {
	if c.domainsURL != "" {
		return defaultDomainsURL
	}
	return c.domainsURL
}

func (c config) GetAllDomains() []string {
	ret := make([]string, 0, len(c.credentials))
	for k := range c.credentials {
		ret = append(ret, k)
	}
	return ret
}

func (c config) GetCredentialsByDomain(domain string) (Credential, error) {
	if v, ok := c.credentials[domain]; ok {
		return v, nil
	}
	return Credential{}, fmt.Errorf("No credentials for %s", domain)
}
