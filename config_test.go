package gdDDNSClient

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_NewConfigFromFile_OK(t *testing.T) {
	confContents := `---
domainsURL: http://domains.google.com
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
		cred, err := conf.GetCredentialsByDomain("test.example.com")
		assert.NoError(err, "GetCredentialsByDomain should not return an error")
		assert.Equal(Credential{"example", "abcd1234"}, cred, "GetCredentialsByDomain should return the expected credentials")
	})
}

func Test_GetDomainsURL_Default(t *testing.T) {
	c := NewConfig("", nil)
	assert.Equal(t, "https://domains.google.com", c.GetDomainsURL())
}

func Test_GetCredentialsByDomain_NoCreds(t *testing.T) {
	creds := map[string]Credential{}
	c := NewConfig("", creds)
	_, err := c.GetCredentialsByDomain("nonexistent")
	assert.Error(t, err)
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
