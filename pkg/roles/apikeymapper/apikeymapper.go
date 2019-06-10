package apikeymapper

import (
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/go-phorce/dolly/xhttp/identity"
	"github.com/go-phorce/dolly/xlog"
	"github.com/go-phorce/dolly/xpki/certutil"
	"github.com/juju/errors"
	yaml "gopkg.in/yaml.v2"
)

var logger = xlog.NewPackageLogger("github.com/certcentral/enrollme/pkg", "apikeymapper")

// ProviderName is identifier for role mapper provider
const ProviderName = "apikey"

const (
	// APIKeyHeader is HTTP header name to be used for auth
	APIKeyHeader = "X-DC-AUTH"
)

// Identity of the caller
type Identity struct {
	// Name of identity
	Name string `json:"name" yaml:"name"`
	// Role of identity
	Role string `json:"role" yaml:"role"`
}

// Config provides mapping of API Keys to Roles
type Config struct {
	// HeaderName is HTTP header
	HeaderName string `json:"header" yaml:"header"`
	// KeysMap is a map of key to Identity
	KeysMap map[string]Identity `json:"keys" yaml:"keys"`
}

// Provider of API Key identity
type Provider struct {
	httpHeader string
	keysMap    map[string]identity.Identity
	namesMap   map[string]identity.Identity
}

// LoadConfig returns configuration loaded from a file
func LoadConfig(file string) (*Config, error) {
	if file == "" {
		return &Config{}, nil
	}

	yamlFile, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, errors.Trace(err)
	}

	var config Config
	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		return nil, errors.Annotatef(err, "unable to unmarshal %q", file)
	}

	if config.HeaderName == "" {
		config.HeaderName = APIKeyHeader
	}

	return &config, nil
}

// Load returns new Provider
func Load(cfgfile string) (*Provider, error) {
	cfg, err := LoadConfig(cfgfile)
	if err != nil {
		return nil, errors.Trace(err)
	}
	return New(cfg), nil
}

// New returns new Provider
func New(cfg *Config) *Provider {
	p := &Provider{
		httpHeader: cfg.HeaderName,
		keysMap:    map[string]identity.Identity{},
		namesMap:   map[string]identity.Identity{},
	}

	if p.httpHeader == "" {
		p.httpHeader = APIKeyHeader
	}

	for key, id := range cfg.KeysMap {
		key = strings.ToUpper(key)
		p.keysMap[key] = identity.NewIdentity(id.Role, id.Name, "")
	}

	return p
}

// HTTPHeaderName returns name of HTTP header to be used in auth
func (p *Provider) HTTPHeaderName() string {
	return p.httpHeader
}

// Applicable returns true if the request has autherization data applicable to the provider
func (p *Provider) Applicable(r *http.Request) bool {
	key := r.Header.Get(p.httpHeader)
	return key != ""
}

// IdentityMapper interface
func (p *Provider) IdentityMapper(r *http.Request) (identity.Identity, error) {
	key := r.Header.Get(p.httpHeader)
	if key == "" {
		return nil, nil
	}
	key = strings.ToUpper(certutil.SHA256Hex([]byte(key)))
	if id, ok := p.keysMap[key]; ok {
		logger.Infof("api=IdentityMapper, role=%s, name=%q", id.Role(), id.Name())
		return id, nil
	}
	return nil, errors.New("invalid access key")
}
