package main

import (
	"errors"
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	client "github.com/twexler/gd-ddns-client"
)

type mockDomainsClient struct {
	calls int
	err   error
}
type mockIPifyClient struct {
	calls int
	err   error
	ip    net.IP
}

func (m *mockDomainsClient) Update(_ client.Credential, _ string, _ net.IP, _ bool) error {
	m.calls = m.calls + 1
	return m.err
}

func (m *mockIPifyClient) GetIPAddress() (net.IP, error) {
	m.calls = m.calls + 1
	return m.ip, m.err
}

func Test_runUpdates_OK(t *testing.T) {
	timeChan := make(chan time.Time)
	creds := client.Credential{}
	ipifyClient := mockIPifyClient{
		ip: net.ParseIP("1.2.3.4"),
	}
	domainsClient := mockDomainsClient{}
	// lol
	currentIP = nil
	go runUpdates(timeChan, "test.example.com", creds, &domainsClient, &ipifyClient)
	timeChan <- time.Now()
	close(timeChan)
	assert.Equal(t, 1, ipifyClient.calls)
	assert.Equal(t, 1, domainsClient.calls)

}

func Test_runUpdates_NoUpdate(t *testing.T) {
	timeChan := make(chan time.Time)
	creds := client.Credential{}
	ipifyClient := mockIPifyClient{}
	domainsClient := mockDomainsClient{}
	// lol
	currentIP = nil
	go runUpdates(timeChan, "test.example.com", creds, &domainsClient, &ipifyClient)
	timeChan <- time.Now()
	close(timeChan)
	assert.Equal(t, 1, ipifyClient.calls)
	assert.Equal(t, 0, domainsClient.calls)
}

func Test_runUpdates_FailUpdate(t *testing.T) {
	timeChan := make(chan time.Time)
	creds := client.Credential{}
	ipifyClient := mockIPifyClient{
		ip: net.ParseIP("1.2.3.4"),
	}
	domainsClient := mockDomainsClient{
		err: errors.New("an error"),
	}
	// lol
	currentIP = nil
	go runUpdates(timeChan, "test.example.com", creds, &domainsClient, &ipifyClient)
	timeChan <- time.Now()
	close(timeChan)
	assert.Equal(t, 1, ipifyClient.calls)
	assert.Equal(t, 1, domainsClient.calls)
}

func Test_runUpdates_FailIPify(t *testing.T) {
	timeChan := make(chan time.Time)
	creds := client.Credential{}
	ipifyClient := mockIPifyClient{
		err: errors.New("an error"),
	}
	domainsClient := mockDomainsClient{}
	go runUpdates(timeChan, "test.example.com", creds, &domainsClient, &ipifyClient)
	timeChan <- time.Now()
	close(timeChan)
	assert.Equal(t, 1, ipifyClient.calls)
	assert.Equal(t, 0, domainsClient.calls)
}
