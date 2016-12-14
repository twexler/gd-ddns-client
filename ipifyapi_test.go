package gdDDNSClient

import (
	"net"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_GetIPAddress_OK(t *testing.T) {
	sm := http.NewServeMux()
	sm.HandleFunc("/", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("1.2.3.4"))
	})
	server := httptest.NewServer(sm)
	defer server.Close()
	c := NewIPIfyAPI(server.URL)
	ip, err := c.GetIPAddress()
	assert.NoError(t, err)
	assert.Equal(t, net.ParseIP("1.2.3.4"), ip)
}
