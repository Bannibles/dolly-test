package roles_test

import (
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"net/http"
	"testing"

	"github.com/certcentral/enrollme/pkg/roles"
	"github.com/certcentral/enrollme/pkg/roles/apikeymapper"
	"github.com/go-phorce/dolly/xhttp/identity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_Empty(t *testing.T) {
	p, err := roles.New("", "", "")
	require.NoError(t, err)

	r, _ := http.NewRequest(http.MethodGet, "/", nil)
	id, err := p.IdentityMapper(r)
	require.NoError(t, err)
	require.NotNil(t, id)
	assert.Equal(t, identity.GuestRoleName, id.Role())
}

func Test_Notfound(t *testing.T) {
	_, err := roles.New("missing_roles.yaml", "", "")
	require.Error(t, err)
	assert.Equal(t, "failed to load JWT mapper: open missing_roles.yaml: no such file or directory", err.Error())

	_, err = roles.New("", "missing_roles.yaml", "")
	require.Error(t, err)
	assert.Equal(t, "failed to load API-Key mapper: open missing_roles.yaml: no such file or directory", err.Error())

	_, err = roles.New("", "", "missing_roles.yaml")
	require.Error(t, err)
	assert.Equal(t, "failed to load cert mapper: open missing_roles.yaml: no such file or directory", err.Error())
}

func Test_All(t *testing.T) {
	p, err := roles.New(
		"jwtmapper/testdata/roles.yaml",
		"apikeymapper/testdata/roles.yaml",
		"certmapper/testdata/roles.yaml")
	require.NoError(t, err)

	t.Run("apikeymapper", func(t *testing.T) {
		r, _ := http.NewRequest(http.MethodGet, "/", nil)
		r.Header.Set(apikeymapper.APIKeyHeader, "A5D461AFA800FBCD820F36707FE9B3051832DB59FD491E6DA91F5E7D884C31C8")

		id, err := p.IdentityMapper(r)
		require.NoError(t, err)
		assert.Equal(t, "admin/Denis Issoupov", id.String())
	})

	t.Run("enrollme-peer", func(t *testing.T) {
		r, _ := http.NewRequest(http.MethodGet, "/", nil)
		r.TLS = &tls.ConnectionState{
			VerifiedChains: [][]*x509.Certificate{
				{
					{
						Subject: pkix.Name{
							CommonName: "[TEST] Digicert Enrollme Root CA",
						},
					},
				},
			},
			PeerCertificates: []*x509.Certificate{
				{
					Subject: pkix.Name{
						CommonName:   "peers.enrollme2.ekspand.com",
						Organization: []string{"Digicert"},
						Country:      []string{"US"},
						Province:     []string{"wa"},
						Locality:     []string{"Kirkland"},
					},
				},
			},
		}
		id, err := p.IdentityMapper(r)
		require.NoError(t, err)
		assert.Equal(t, "enrollme-peer/peers.enrollme2.ekspand.com", id.String())
	})
}
