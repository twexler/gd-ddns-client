package gdDDNSClient

import (
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"time"
)

const (
	defaultIPIfyURL     = "https://api.ipify.org"
	defaultIPIfyTimeout = time.Second
)

// IPIfyAPI is a simple implementation of the ipify.org API
type IPIfyAPI interface {
	// GetIPAddress returns a net.IP or error based on a response from ipify.org's API
	GetIPAddress() (net.IP, error)
}

type ipifyAPI struct {
	client  http.Client
	baseURL string
}

// NewIPIfyAPI returns an instance of an IPIfyAPI
func NewIPIfyAPI(baseURL string) IPIfyAPI {
	client := http.Client{
		Timeout: defaultIPIfyTimeout,
	}
	return &ipifyAPI{
		client:  client,
		baseURL: baseURL,
	}
}

func (i ipifyAPI) GetIPAddress() (net.IP, error) {
	url, err := url.Parse(i.baseURL)
	if err != nil {
		return nil, err
	}
	resp, err := i.client.Get(url.String())
	defer resp.Body.Close()
	if err != nil {
		return nil, err
	}
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return net.ParseIP(string(data)), nil
}
