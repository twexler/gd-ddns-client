package gdDDNSClient

import (
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_NewConfigFromFile_OK(t *testing.T) {
	confContents := `---
domainsURL: http://domains.google.com
ipifyURL: http://api.ipify.org
updateInterval: 1s
credentials:
    test.example.com:
        user: example
        password: abcd1234
`
	withTempFileAndContents(confContents, func(filename string) {
		conf, err := NewConfigFromFile(filename)
		require.NoError(t, err, "NewConfigFromFile should not return an error when loading a valid config")
		assert := assert.New(t)
		assert.NotNil(conf, "configuration should not be nil")
		assert.Equal("http://domains.google.com", conf.GetDomainsURL(), "GetDomainsURL should return the expected string")
		assert.Len(conf.GetAllDomains(), 1, "GetAllDomains should return a slice of length 1")
		cred, err := conf.GetCredentialsByHostname("test.example.com")
		assert.NoError(err, "GetCredentialsByDomain should not return an error")
		assert.Equal(Credential{"example", "abcd1234"}, cred, "GetCredentialsByDomain should return the expected credentials")
		assert.Equal(time.Second, conf.GetUpdateInterval(), "Interval should equal 1s")
		assert.Equal("http://api.ipify.org", conf.GetIPIfyURL(), "IPify API url should be equal to the config value")
	})
}

func Test_Defaults(t *testing.T) {
	creds := map[string]Credential{}
	c := config{Credentials: creds}
	_, err := c.GetCredentialsByHostname("nonexistent")
	assert := assert.New(t)
	assert.Error(err)
	assert.Equal(defaultUpdateInterval, c.GetUpdateInterval())
	assert.Equal(defaultDomainsURL, c.GetDomainsURL())
	assert.Equal(defaultIPIfyURL, c.GetIPIfyURL())
}

func withTempFileAndContents(contents string, fn func(filename string)) {
	f, err := ioutil.TempFile("/tmp", "gd-ddns-testing")
	if err != nil {
		panic(err)
	}
	defer os.Remove(f.Name())
	_, err = f.WriteString(contents)
	if err != nil {
		panic(err)
	}
	fn(f.Name())
}
