package gdDDNSClient

import (
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	defaultDomainsAPITimeout = time.Second
	goodUpdateResponse       = "good"
	noChangeUpdateResponse   = "nochg"
	updateURI                = "/nic/update"
)

// DomainsAPI is an interface that Domains API implementations must implement
type DomainsAPI interface {
	// Update conforms to the HTTP API for updating a DDNS name on Google Domains
	Update(credential Credential, hostname string, ip net.IP, offline bool) error
}

type domainsAPI struct {
	client  http.Client
	baseURL string
}

// NewDomainsAPI returns a new instance of a DaomainsAPI
func NewDomainsAPI(baseURL string) DomainsAPI {
	client := http.Client{
		Timeout: defaultDomainsAPITimeout,
	}
	return &domainsAPI{
		client:  client,
		baseURL: baseURL,
	}
}

func (a domainsAPI) Update(credential Credential, hostname string, ip net.IP, offline bool) error {
	queryValues := make(url.Values, 3) // query values are hostname, ip and offline
	queryValues.Add("hostname", hostname)
	queryValues.Add("ip", ip.String())
	// offline can be "yes" or "no"
	offlineValue := "no"
	if offline {
		offlineValue = "yes"
	}
	queryValues.Add("offline", offlineValue)

	url, err := url.Parse(a.baseURL)
	if err != nil {
		return err
	}
	url.Path = updateURI
	url.RawQuery = queryValues.Encode()

	req, err := http.NewRequest("GET", url.String(), nil)
	if err != nil {
		return err
	}
	req.SetBasicAuth(credential.User, credential.Password)
	resp, err := a.client.Do(req)
	if err != nil {
		return err
	}
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	resp.Body.Close()

	// stringify things
	respString := string(data)
	if strings.HasPrefix(respString, goodUpdateResponse) || strings.HasPrefix(respString, noChangeUpdateResponse) {
		return nil
	}

	return fmt.Errorf("Got error from Google Domains: %s", respString)
}
