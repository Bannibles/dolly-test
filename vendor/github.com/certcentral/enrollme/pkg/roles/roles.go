package roles

import (
	"net/http"

	"github.com/certcentral/enrollme/pkg/roles/apikeymapper"
	"github.com/certcentral/enrollme/pkg/roles/certmapper"
	"github.com/certcentral/enrollme/pkg/roles/jwtmapper"
	"github.com/go-phorce/dolly/xhttp/identity"
	"github.com/go-phorce/dolly/xlog"
	"github.com/juju/errors"
)

var logger = xlog.NewPackageLogger("github.com/certcentral/enrollme/pkg", "roles")

// IdentityProvider interface to extract identity from requests
type IdentityProvider interface {
	// Applicable returns true if the provider is applicable for the request
	Applicable(*http.Request) bool
	// IdentityMapper returns identity from the request
	IdentityMapper(*http.Request) (identity.Identity, error)
}

// Provider for authz identity
type Provider struct {
	CertMapper   *certmapper.Provider
	JwtMapper    *jwtmapper.Provider
	APIkeyMapper *apikeymapper.Provider
}

// New returns Authz provider instance
func New(jwtMapper, apiKeyMapper, certMapper string) (*Provider, error) {
	var err error
	prov := new(Provider)

	if certMapper != "" {
		prov.CertMapper, err = certmapper.Load(certMapper)
		if err != nil {
			return nil, errors.Annotatef(err, "failed to load cert mapper")
		}
	}
	if jwtMapper != "" {
		prov.JwtMapper, err = jwtmapper.Load(jwtMapper)
		if err != nil {
			return nil, errors.Annotatef(err, "failed to load JWT mapper")
		}
	}
	if apiKeyMapper != "" {
		prov.APIkeyMapper, err = apikeymapper.Load(apiKeyMapper)
		if err != nil {
			return nil, errors.Annotatef(err, "failed to load API-Key mapper")
		}
	}

	return prov, nil
}

// IdentityMapper returns identity from the request
func (p *Provider) IdentityMapper(r *http.Request) (identity.Identity, error) {
	if p.JwtMapper != nil && p.JwtMapper.Applicable(r) {
		return p.JwtMapper.IdentityMapper(r)
	}
	if p.APIkeyMapper != nil && p.APIkeyMapper.Applicable(r) {
		return p.APIkeyMapper.IdentityMapper(r)
	}
	if p.CertMapper != nil && p.CertMapper.Applicable(r) {
		return p.CertMapper.IdentityMapper(r)
	}

	// if none of mappers are applicable or configured,
	// then use default guest mapper
	return identity.GuestIdentityMapper(r)
}
