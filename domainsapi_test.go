package gdDDNSClient

import (
	"net"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_Update_OK(t *testing.T) {
	withMockDomainsAPI(t, []byte("good 1.2.3.4"), func(a DomainsAPI) {
		cred := Credential{
			User:     "myuser",
			Password: "mypassword",
		}
		ip := net.ParseIP("1.2.3.4")
		assert.NoError(t, a.Update(cred, "foo.bar.baz", ip, false))
	})
}

func withMockDomainsAPI(t *testing.T, response []byte, fn func(DomainsAPI)) {
	sm := http.NewServeMux()
	sm.HandleFunc("/nic/update", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write(response)
	})
	server := httptest.NewServer(sm)
	defer server.Close()
	c := NewDomainsAPI(server.URL)
	require.NotNil(t, c)
	fn(c)
}
