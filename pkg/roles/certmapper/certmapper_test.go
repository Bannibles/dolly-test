package certmapper_test

import (
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"net/http"
	"testing"

	"github.com/certcentral/enrollme/pkg/roles/certmapper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_Config(t *testing.T) {
	_, err := certmapper.Load("testdata/missing.yaml")
	require.Error(t, err)
	assert.Equal(t, "open testdata/missing.yaml: no such file or directory", err.Error())

	_, err = certmapper.Load("testdata/roles_corrupted.1.yaml")
	require.Error(t, err)
	assert.Equal(t, `unable to unmarshal "testdata/roles_corrupted.1.yaml": yaml: line 2: mapping values are not allowed in this context`, err.Error())

	_, err = certmapper.Load("testdata/roles_corrupted.2.yaml")
	require.Error(t, err)
	assert.Equal(t, `unable to unmarshal "testdata/roles_corrupted.2.yaml": yaml: line 4: could not find expected ':'`, err.Error())

	_, err = certmapper.Load("")
	require.NoError(t, err)

	cfg, err := certmapper.LoadConfig("testdata/roles.yaml")
	require.NoError(t, err)
	assert.Equal(t, 2, len(cfg.ValidOrganizations))
	assert.Equal(t, 1, len(cfg.ValidIssuers))
}

func Test_identity(t *testing.T) {
	cfg, err := certmapper.LoadConfig("testdata/roles.yaml")
	require.NoError(t, err)

	p := certmapper.New(cfg)

	t.Run("not_applicable", func(t *testing.T) {
		r, _ := http.NewRequest(http.MethodGet, "/", nil)
		id, err := p.IdentityMapper(r)
		require.NoError(t, err)
		assert.Nil(t, id)
	})

	t.Run("org_not_allowed", func(t *testing.T) {
		r, _ := http.NewRequest(http.MethodGet, "/", nil)
		r.TLS = &tls.ConnectionState{
			PeerCertificates: []*x509.Certificate{
				{
					Subject: pkix.Name{
						CommonName:   "dolly",
						Organization: []string{"org"},
					},
				},
			},
		}
		_, err := p.IdentityMapper(r)
		require.Error(t, err)
		assert.Equal(t, `the "org" organization is not allowed`, err.Error())
	})

	t.Run("no_issuer", func(t *testing.T) {
		r, _ := http.NewRequest(http.MethodGet, "/", nil)
		r.TLS = &tls.ConnectionState{
			PeerCertificates: []*x509.Certificate{
				{
					Subject: pkix.Name{
						CommonName:   "dolly",
						Organization: []string{"Digicert"},
					},
				},
			},
		}
		_, err := p.IdentityMapper(r)
		require.Error(t, err)
		assert.Equal(t, `the "" root CA is not allowed`, err.Error())
	})

	t.Run("not_issuer", func(t *testing.T) {
		r, _ := http.NewRequest(http.MethodGet, "/", nil)
		r.TLS = &tls.ConnectionState{
			VerifiedChains: [][]*x509.Certificate{
				{
					{
						Subject: pkix.Name{
							CommonName: "issuer",
						},
					},
				},
			},
			PeerCertificates: []*x509.Certificate{
				{
					Subject: pkix.Name{
						CommonName:   "dolly",
						Organization: []string{"Digicert"},
					},
				},
			},
		}
		_, err := p.IdentityMapper(r)
		require.Error(t, err)
		assert.Equal(t, `the "CN=issuer" root CA is not allowed`, err.Error())
	})

	t.Run("not_in_map_CN", func(t *testing.T) {
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
						CommonName:   "dolly",
						Organization: []string{"Digicert"},
					},
				},
			},
		}
		_, err := p.IdentityMapper(r)
		require.Error(t, err)
		assert.Equal(t, `the "O=Digicert, CN=dolly" subject is not allowed`, err.Error())
	})

	t.Run("enrollme-client", func(t *testing.T) {
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
						CommonName:   "enrollme1.ekspand.com",
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
		assert.Equal(t, "enrollme-client/enrollme1.ekspand.com", id.String())
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
