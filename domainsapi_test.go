package gdDDNSClient

import (
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_Update_OK(t *testing.T) {
	for _, v := range []bool{true, false} {
		v := v
		t.Run(fmt.Sprintf("offline %t"), func(st *testing.T) {
			withMockDomainsAPI(st, []byte("good 1.2.3.4"), func(a DomainsAPI) {
				cred := Credential{
					User:     "myuser",
					Password: "mypassword",
				}
				ip := net.ParseIP("1.2.3.4")
				assert.NoError(st, a.Update(cred, "foo.bar.baz", ip, v))
			})
		})
	}
}

func Test_Update_Fail_ErrorResponse(t *testing.T) {
	withMockDomainsAPI(t, []byte("badauth"), func(a DomainsAPI) {
		cred := Credential{
			User:     "myuser",
			Password: "mypassword",
		}
		ip := net.ParseIP("1.2.3.4")
		assert.Error(t, a.Update(cred, "foo.bar.baz", ip, false))
	})
}

func Test_Update_Fail_InvalidAPIURL(t *testing.T) {
	c := NewDomainsAPI("invalid:/:///%%%%")
	assert.Error(t, c.Update(Credential{}, "", nil, true))
}

func Test_Update_Fail_ConnectionFailure(t *testing.T) {
	c := NewDomainsAPI("http://1.2.3.4:9999")
	assert.Error(t, c.Update(Credential{}, "", nil, true))
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
