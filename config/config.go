// Package config allows for an external config file to be read that allows for
// value to be overriden based on a hostname derived configuration set.
//
// the Configuration type defines all the configurable parameters.
// the config file is json, its consists of 3 sections
//
// defaults   : a Configuration instance that is the base/default configurations
// hosts      : a mapping from host name to a named configuration [e.g. node1 : "aws"]
// overrrides : a set of named Configuration instances that can override the some or all of the default config values
//
// the caller can provide a specific hostname if it chooses, otherwise the config will
//  a) look for a named environemnt variable, if set to something, that is used
//  b) look at the OS supplied hostname
//
//
// *** THIS IS GENERATED CODE: DO NOT EDIT ***
package config

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"
)

// Duration represents a period of time, its the same as time.Duration
// but supports better marshalling from json
type Duration time.Duration

// UnmarshalJSON handles decoding our custom json serialization for Durations
// json values that are numbers are treated as seconds
// json values that are strings, can use the standard time.Duration units indicators
// e.g. this can decode val:100 as well as val:"10m"
func (d *Duration) UnmarshalJSON(b []byte) error {
	if b[0] == '"' {
		dir, err := time.ParseDuration(string(b[1 : len(b)-1]))
		*d = Duration(dir)
		return err
	}
	i, err := json.Number(string(b)).Int64()
	*d = Duration(time.Duration(i) * time.Second)
	return err
}

// MarshalJSON encodes our custom Duration value as a quoted version of its underlying value's String() output
// this means you get a duration with a trailing units indicator, e.g. "10m0s"
func (d Duration) MarshalJSON() ([]byte, error) {
	return []byte(`"` + d.String() + `"`), nil
}

// String returns a string formatted version of the duration in a valueUnits format, e.g. 5m0s for 5 minutes
func (d Duration) String() string {
	return time.Duration(d).String()
}

// TimeDuration returns this duration in a time.Duration type
func (d Duration) TimeDuration() time.Duration {
	return time.Duration(d)
}

// Authz contains configuration for the authorization module.
type Authz struct {

	// Allow will allow the specified roles access to this path and its children, in format: ${path}:${role},${role}.
	Allow []string

	// AllowAny will allow any authenticated request access to this path and its children.
	AllowAny []string

	// AllowAnyRole will allow any authenticated request that include a non empty role.
	AllowAnyRole []string

	// LogAllowed specifies to log allowed access.
	LogAllowed *bool

	// LogDenied specifies to log denied access.
	LogDenied *bool

	// CertMapper specifies location of the config file for certificate based identity.
	CertMapper string

	// APIKeyMapper specifies location of the config file for API-Key based identity.
	APIKeyMapper string

	// JWTMapper specifies location of the config file for JWT based identity.
	JWTMapper string
}

func (c *Authz) overrideFrom(o *Authz) {
	overrideStrings(&c.Allow, &o.Allow)
	overrideStrings(&c.AllowAny, &o.AllowAny)
	overrideStrings(&c.AllowAnyRole, &o.AllowAnyRole)
	overrideBool(&c.LogAllowed, &o.LogAllowed)
	overrideBool(&c.LogDenied, &o.LogDenied)
	overrideString(&c.CertMapper, &o.CertMapper)
	overrideString(&c.APIKeyMapper, &o.APIKeyMapper)
	overrideString(&c.JWTMapper, &o.JWTMapper)

}

// AuthzConfig contains configuration for the authorization module.
type AuthzConfig interface {
	// Allow will allow the specified roles access to this path and its children, in format: ${path}:${role},${role}.
	GetAllow() []string
	// AllowAny will allow any authenticated request access to this path and its children.
	GetAllowAny() []string
	// AllowAnyRole will allow any authenticated request that include a non empty role.
	GetAllowAnyRole() []string
	// LogAllowed specifies to log allowed access.
	GetLogAllowed() bool
	// LogDenied specifies to log denied access.
	GetLogDenied() bool
	// CertMapper specifies location of the config file for certificate based identity.
	GetCertMapper() string
	// APIKeyMapper specifies location of the config file for API-Key based identity.
	GetAPIKeyMapper() string
	// JWTMapper specifies location of the config file for JWT based identity.
	GetJWTMapper() string
}

// GetAllow will allow the specified roles access to this path and its children, in format: ${path}:${role},${role}.
func (c *Authz) GetAllow() []string {
	return c.Allow
}

// GetAllowAny will allow any authenticated request access to this path and its children.
func (c *Authz) GetAllowAny() []string {
	return c.AllowAny
}

// GetAllowAnyRole will allow any authenticated request that include a non empty role.
func (c *Authz) GetAllowAnyRole() []string {
	return c.AllowAnyRole
}

// GetLogAllowed specifies to log allowed access.
func (c *Authz) GetLogAllowed() bool {
	return c.LogAllowed != nil && *c.LogAllowed
}

// GetLogDenied specifies to log denied access.
func (c *Authz) GetLogDenied() bool {
	return c.LogDenied != nil && *c.LogDenied
}

// GetCertMapper specifies location of the config file for certificate based identity.
func (c *Authz) GetCertMapper() string {
	return c.CertMapper
}

// GetAPIKeyMapper specifies location of the config file for API-Key based identity.
func (c *Authz) GetAPIKeyMapper() string {
	return c.APIKeyMapper
}

// GetJWTMapper specifies location of the config file for JWT based identity.
func (c *Authz) GetJWTMapper() string {
	return c.JWTMapper
}

// CORS contains configuration for CORS.
type CORS struct {

	// Enabled specifies if the CORS is enabled.
	Enabled *bool

	// MaxAge indicates how long (in seconds) the results of a preflight request can be cached.
	MaxAge int

	// AllowedOrigins is a list of origins a cross-domain request can be executed from.
	AllowedOrigins []string

	// AllowedMethods is a list of methods the client is allowed to use with cross-domain requests.
	AllowedMethods []string

	// AllowedHeaders is list of non simple headers the client is allowed to use with cross-domain requests.
	AllowedHeaders []string

	// ExposedHeaders indicates which headers are safe to expose to the API of a CORS API specification.
	ExposedHeaders []string

	// AllowCredentials indicates whether the request can include user credentials.
	AllowCredentials *bool

	// OptionsPassthrough instructs preflight to let other potential next handlers to process the OPTIONS method.
	OptionsPassthrough *bool

	// Debug flag adds additional output to debug server side CORS issues.
	Debug *bool
}

func (c *CORS) overrideFrom(o *CORS) {
	overrideBool(&c.Enabled, &o.Enabled)
	overrideInt(&c.MaxAge, &o.MaxAge)
	overrideStrings(&c.AllowedOrigins, &o.AllowedOrigins)
	overrideStrings(&c.AllowedMethods, &o.AllowedMethods)
	overrideStrings(&c.AllowedHeaders, &o.AllowedHeaders)
	overrideStrings(&c.ExposedHeaders, &o.ExposedHeaders)
	overrideBool(&c.AllowCredentials, &o.AllowCredentials)
	overrideBool(&c.OptionsPassthrough, &o.OptionsPassthrough)
	overrideBool(&c.Debug, &o.Debug)

}

// CORSConfig contains configuration for CORSConfig.
type CORSConfig interface {
	// Enabled specifies if the CORS is enabled.
	GetEnabled() bool
	// MaxAge indicates how long (in seconds) the results of a preflight request can be cached.
	GetMaxAge() int
	// AllowedOrigins is a list of origins a cross-domain request can be executed from.
	GetAllowedOrigins() []string
	// AllowedMethods is a list of methods the client is allowed to use with cross-domain requests.
	GetAllowedMethods() []string
	// AllowedHeaders is list of non simple headers the client is allowed to use with cross-domain requests.
	GetAllowedHeaders() []string
	// ExposedHeaders indicates which headers are safe to expose to the API of a CORS API specification.
	GetExposedHeaders() []string
	// AllowCredentials indicates whether the request can include user credentials.
	GetAllowCredentials() bool
	// OptionsPassthrough instructs preflight to let other potential next handlers to process the OPTIONS method.
	GetOptionsPassthrough() bool
	// Debug flag adds additional output to debug server side CORS issues.
	GetDebug() bool
}

// GetEnabled specifies if the CORS is enabled.
func (c *CORS) GetEnabled() bool {
	return c.Enabled != nil && *c.Enabled
}

// GetMaxAge indicates how long (in seconds) the results of a preflight request can be cached.
func (c *CORS) GetMaxAge() int {
	return c.MaxAge
}

// GetAllowedOrigins is a list of origins a cross-domain request can be executed from.
func (c *CORS) GetAllowedOrigins() []string {
	return c.AllowedOrigins
}

// GetAllowedMethods is a list of methods the client is allowed to use with cross-domain requests.
func (c *CORS) GetAllowedMethods() []string {
	return c.AllowedMethods
}

// GetAllowedHeaders is list of non simple headers the client is allowed to use with cross-domain requests.
func (c *CORS) GetAllowedHeaders() []string {
	return c.AllowedHeaders
}

// GetExposedHeaders indicates which headers are safe to expose to the API of a CORS API specification.
func (c *CORS) GetExposedHeaders() []string {
	return c.ExposedHeaders
}

// GetAllowCredentials indicates whether the request can include user credentials.
func (c *CORS) GetAllowCredentials() bool {
	return c.AllowCredentials != nil && *c.AllowCredentials
}

// GetOptionsPassthrough instructs preflight to let other potential next handlers to process the OPTIONS method.
func (c *CORS) GetOptionsPassthrough() bool {
	return c.OptionsPassthrough != nil && *c.OptionsPassthrough
}

// GetDebug flag adds additional output to debug server side CORS issues.
func (c *CORS) GetDebug() bool {
	return c.Debug != nil && *c.Debug
}

// Configuration contains the configuration for the server.
type Configuration struct {

	// Datacenter specifies the Datacenter where the instance is running.
	Datacenter string

	// Environment specifies the deployment environment.
	Environment string

	// ServiceName specifies the service name to be used in logs and folders names.
	ServiceName string

	// HTTP contains the config for the Public HTTP.
	HTTP HTTPServer

	// HTTPS contains the config for the HTTPS/JSON API Service.
	HTTPS HTTPServer

	// Authz contains configuration for the API authorization layer.
	Authz Authz

	// Audit contains configuration for the audit logger.
	Audit Logger

	// CryptoProv specifies the configuration for crypto providers.
	CryptoProv CryptoProv

	// Metrics specifies the metrics pipeline configuration.
	Metrics Metrics

	// Logger contains configuration for the logger.
	Logger Logger

	// LogLevels specifies the log levels per package.
	LogLevels []RepoLogLevel

	// RootCA specifies the location of PEM-encoded certificate.
	RootCA string
}

func (c *Configuration) overrideFrom(o *Configuration) {
	overrideString(&c.Datacenter, &o.Datacenter)
	overrideString(&c.Environment, &o.Environment)
	overrideString(&c.ServiceName, &o.ServiceName)
	c.HTTP.overrideFrom(&o.HTTP)
	c.HTTPS.overrideFrom(&o.HTTPS)
	c.Authz.overrideFrom(&o.Authz)
	c.Audit.overrideFrom(&o.Audit)
	c.CryptoProv.overrideFrom(&o.CryptoProv)
	c.Metrics.overrideFrom(&o.Metrics)
	c.Logger.overrideFrom(&o.Logger)
	overrideRepoLogLevelSlice(&c.LogLevels, &o.LogLevels)
	overrideString(&c.RootCA, &o.RootCA)

}

// CryptoProv specifies the configuration for crypto providers.
type CryptoProv struct {

	// Default specifies the location of the configuration file for default provider.
	Default string

	// Providers specifies the list of locations of the configuration files.
	Providers []string
}

func (c *CryptoProv) overrideFrom(o *CryptoProv) {
	overrideString(&c.Default, &o.Default)
	overrideStrings(&c.Providers, &o.Providers)

}

// HTTPServer contains the configuration of the HTTPS API Service.
type HTTPServer struct {

	// ServiceName specifies name of the service: HTTP|HTTPS|WebAPI.
	ServiceName string

	// Disabled specifies if the service is disabled.
	Disabled *bool

	// VIPName is the FQ name of the VIP to the cluster [this is used when building the cert requests].
	VIPName string

	// BindAddr is the address that the HTTPS service should be exposed on
	BindAddr string

	// ServerTLS provides TLS config for server.
	ServerTLS TLSInfo

	// PackageLogger if set, specifies name of the package logger.
	PackageLogger string

	// AllowProfiling if set, will allow for per request CPU/Memory profiling triggered by the URI QueryString.
	AllowProfiling *bool

	// ProfilerDir specifies the directories where per-request profile information is written, if not set will write to a TMP dir.
	ProfilerDir string

	// Services is a list of services to enable for this HTTP Service.
	Services []string

	// HeartbeatSecs specifies heartbeat interval in seconds [30 secs is a minimum].
	HeartbeatSecs int

	// CORS contains configuration for CORS.
	CORS CORS
}

func (c *HTTPServer) overrideFrom(o *HTTPServer) {
	overrideString(&c.ServiceName, &o.ServiceName)
	overrideBool(&c.Disabled, &o.Disabled)
	overrideString(&c.VIPName, &o.VIPName)
	overrideString(&c.BindAddr, &o.BindAddr)
	c.ServerTLS.overrideFrom(&o.ServerTLS)
	overrideString(&c.PackageLogger, &o.PackageLogger)
	overrideBool(&c.AllowProfiling, &o.AllowProfiling)
	overrideString(&c.ProfilerDir, &o.ProfilerDir)
	overrideStrings(&c.Services, &o.Services)
	overrideInt(&c.HeartbeatSecs, &o.HeartbeatSecs)
	c.CORS.overrideFrom(&o.CORS)

}

// HTTPServerConfig contains the configuration of the HTTPS API Service.
type HTTPServerConfig interface {
	// ServiceName specifies name of the service: HTTP|HTTPS|WebAPI.
	GetServiceName() string
	// Disabled specifies if the service is disabled.
	GetDisabled() bool
	// VIPName is the FQ name of the VIP to the cluster [this is used when building the cert requests].
	GetVIPName() string
	// BindAddr is the address that the HTTPS service should be exposed on
	GetBindAddr() string
	// ServerTLS provides TLS config for server.
	GetServerTLSCfg() *TLSInfo
	// PackageLogger if set, specifies name of the package logger.
	GetPackageLogger() string
	// AllowProfiling if set, will allow for per request CPU/Memory profiling triggered by the URI QueryString.
	GetAllowProfiling() bool
	// ProfilerDir specifies the directories where per-request profile information is written, if not set will write to a TMP dir.
	GetProfilerDir() string
	// Services is a list of services to enable for this HTTP Service.
	GetServices() []string
	// HeartbeatSecs specifies heartbeat GetHeartbeatSecserval in seconds [30 secs is a minimum].
	GetHeartbeatSecs() int
	// GetCORSCfg contains configuration for GetCORSCfg.
	GetCORSCfg() *CORS
}

// GetServiceName specifies name of the service: HTTP|HTTPS|WebAPI.
func (c *HTTPServer) GetServiceName() string {
	return c.ServiceName
}

// GetDisabled specifies if the service is disabled.
func (c *HTTPServer) GetDisabled() bool {
	return c.Disabled != nil && *c.Disabled
}

// GetVIPName is the FQ name of the VIP to the cluster [this is used when building the cert requests].
func (c *HTTPServer) GetVIPName() string {
	return c.VIPName
}

// GetBindAddr is the address that the HTTPS service should be exposed on
func (c *HTTPServer) GetBindAddr() string {
	return c.BindAddr
}

// GetServerTLSCfg provides TLS config for server.
func (c *HTTPServer) GetServerTLSCfg() *TLSInfo {
	return &c.ServerTLS
}

// GetPackageLogger if set, specifies name of the package logger.
func (c *HTTPServer) GetPackageLogger() string {
	return c.PackageLogger
}

// GetAllowProfiling if set, will allow for per request CPU/Memory profiling triggered by the URI QueryString.
func (c *HTTPServer) GetAllowProfiling() bool {
	return c.AllowProfiling != nil && *c.AllowProfiling
}

// GetProfilerDir specifies the directories where per-request profile information is written, if not set will write to a TMP dir.
func (c *HTTPServer) GetProfilerDir() string {
	return c.ProfilerDir
}

// GetServices is a list of services to enable for this HTTP Service.
func (c *HTTPServer) GetServices() []string {
	return c.Services
}

// GetHeartbeatSecs specifies heartbeat interval in seconds [30 secs is a minimum].
func (c *HTTPServer) GetHeartbeatSecs() int {
	return c.HeartbeatSecs
}

// GetCORSCfg contains configuration for GetCORSCfg.
func (c *HTTPServer) GetCORSCfg() *CORS {
	return &c.CORS
}

// Logger contains information about the configuration of a logger/log rotation.
type Logger struct {

	// Directory contains where to store the log files.
	Directory string

	// MaxAgeDays controls how old files are before deletion.
	MaxAgeDays int

	// MaxSizeMb contols how large a single log file can be before its rotated.
	MaxSizeMb int
}

func (c *Logger) overrideFrom(o *Logger) {
	overrideString(&c.Directory, &o.Directory)
	overrideInt(&c.MaxAgeDays, &o.MaxAgeDays)
	overrideInt(&c.MaxSizeMb, &o.MaxSizeMb)

}

// LoggerConfig contains information about the configuration of a logger/log rotation.
type LoggerConfig interface {
	// Directory contains where to store the log files.
	GetDirectory() string
	// MaxAgeDays controls how old files are before deletion.
	GetMaxAgeDays() int
	// MaxSizeMb contols how large a single log file can be before its rotated.
	GetMaxSizeMb() int
}

// GetDirectory contains where to store the log files.
func (c *Logger) GetDirectory() string {
	return c.Directory
}

// GetMaxAgeDays controls how old files are before deletion.
func (c *Logger) GetMaxAgeDays() int {
	return c.MaxAgeDays
}

// GetMaxSizeMb contols how large a single log file can be before its rotated.
func (c *Logger) GetMaxSizeMb() int {
	return c.MaxSizeMb
}

// Metrics specifies the metrics pipeline configuration.
type Metrics struct {

	// Provider specifies the metrics provider.
	Provider string
}

func (c *Metrics) overrideFrom(o *Metrics) {
	overrideString(&c.Provider, &o.Provider)

}

// RepoLogLevel contains information about the log level per repo. Use * to set up global level.
type RepoLogLevel struct {

	// Repo specifies the repo name, or '*' for all repos [Global].
	Repo string

	// Package specifies the package name.
	Package string

	// Level specifies the log level for the repo [ERROR,WARNING,NOTICE,INFO,DEBUG,TRACE].
	Level string
}

func (c *RepoLogLevel) overrideFrom(o *RepoLogLevel) {
	overrideString(&c.Repo, &o.Repo)
	overrideString(&c.Package, &o.Package)
	overrideString(&c.Level, &o.Level)

}

// TLSInfo contains configuration info for the TLS.
type TLSInfo struct {

	// CertFile specifies location of the cert.
	CertFile string

	// KeyFile specifies location of the key.
	KeyFile string

	// TrustedCAFile specifies location of the trusted Root file.
	TrustedCAFile string

	// ClientCertAuth controls client auth: NoClientCert|RequestClientCert|RequireAnyClientCert|VerifyClientCertIfGiven|RequireAndVerifyClientCert
	ClientCertAuth string
}

func (c *TLSInfo) overrideFrom(o *TLSInfo) {
	overrideString(&c.CertFile, &o.CertFile)
	overrideString(&c.KeyFile, &o.KeyFile)
	overrideString(&c.TrustedCAFile, &o.TrustedCAFile)
	overrideString(&c.ClientCertAuth, &o.ClientCertAuth)

}

// TLSInfoConfig contains configuration info for the TLS.
type TLSInfoConfig interface {
	// CertFile specifies location of the cert.
	GetCertFile() string
	// KeyFile specifies location of the key.
	GetKeyFile() string
	// TrustedCAFile specifies location of the trusted Root file.
	GetTrustedCAFile() string
	// ClientCertAuth controls client auth: NoClientCert|RequestClientCert|RequireAnyClientCert|VerifyClientCertIfGiven|RequireAndVerifyClientCert
	GetClientCertAuth() string
}

// GetCertFile specifies location of the cert.
func (c *TLSInfo) GetCertFile() string {
	return c.CertFile
}

// GetKeyFile specifies location of the key.
func (c *TLSInfo) GetKeyFile() string {
	return c.KeyFile
}

// GetTrustedCAFile specifies location of the trusted Root file.
func (c *TLSInfo) GetTrustedCAFile() string {
	return c.TrustedCAFile
}

// GetClientCertAuth controls client auth: NoClientCert|RequestClientCert|RequireAnyClientCert|VerifyClientCertIfGiven|RequireAndVerifyClientCert
func (c *TLSInfo) GetClientCertAuth() string {
	return c.ClientCertAuth
}

func overrideBool(d, o **bool) {
	if *o != nil {
		*d = *o
	}
}

func overrideInt(d, o *int) {
	if *o != 0 {
		*d = *o
	}
}

func overrideRepoLogLevelSlice(d, o *[]RepoLogLevel) {
	if len(*o) > 0 {
		*d = *o
	}
}

func overrideString(d, o *string) {
	if *o != "" {
		*d = *o
	}
}

func overrideStrings(d, o *[]string) {
	if len(*o) > 0 {
		*d = *o
	}
}

// Load will attempt to load the configuration from the supplied filename.
// Overrides defined in the config file will be applied based on the hostname
// the hostname used is dervied from [in order]
//    1) the hostnameOverride parameter if not ""
//    2) the value of the Environment variable in envKeyName, if not ""
//    3) the OS supplied hostname
func Load(configFilename, envKeyName, hostnameOverride string) (*Configuration, error) {
	configs, err := LoadConfigurations(configFilename)
	if err != nil {
		return nil, err
	}
	return configs.For(envKeyName, hostnameOverride)
}

// LoadConfigurations decodes the json config file, or returns an error
// typically you'd just use Load, but this can be useful if you need to
// do more intricate examination of the entire set of configurations
func LoadConfigurations(filename string) (*Configurations, error) {
	cfr, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer cfr.Close()
	configs := new(Configurations)
	return configs, json.NewDecoder(cfr).Decode(configs)
}

// Configurations is the entire set of configurations, these consist of
//    a base/default configuration
//    a set of hostname -> named overrides
//    named overrides -> config overrides
type Configurations struct {
	// Default contains the base configuration, this applies unless it override by a specifc named config
	Defaults Configuration

	// a map of hostname to named configuration
	Hosts map[string]string

	// a map of named configuration overrides
	Overrides map[string]Configuration
}

// HostSelection describes the hostname & override set that were used
type HostSelection struct {
	// Hostname returns the hostname from the configuration that was used
	// this may return a fully qualified hostname, when just a name was specified
	Hostname string
	// Override contains the name of the override section, if there was one found
	// [based on the Hostname]
	Override string
}

// For returns the Configuration for the indicated host, with all the overrides applied.
// the hostname used is dervied from [in order]
//    1) the hostnameOverride parameter if not ""
//    2) the value of the Environemnt variable in envKeyName, if not ""
//    3) the OS supplied hostname
func (configs *Configurations) For(envKeyName, hostnameOverride string) (*Configuration, error) {
	sel, err := configs.Selection(envKeyName, hostnameOverride)
	if err != nil {
		return nil, err
	}
	c := configs.Defaults
	if sel.Override != "" {
		overrides := configs.Overrides[sel.Override]
		c.overrideFrom(&overrides)
	}
	return &c, nil
}

// Selection returns the final resolved hostname, and if applicable,
// override section name for the supplied host specifiers
func (configs *Configurations) Selection(envKeyName, hostnameOverride string) (HostSelection, error) {
	res := HostSelection{}
	hn, err := configs.resolveHostname(envKeyName, hostnameOverride)
	if err != nil {
		return res, err
	}
	res.Hostname = hn
	if ov, exists := configs.Hosts[hn]; exists {
		if _, exists := configs.Overrides[ov]; !exists {
			return res, fmt.Errorf("Configuration for host %s specified override set %s but that doesn't exist", hn, ov)
		}
		res.Override = ov
	}
	return res, nil
}

// resolveHostname determines the hostname to lookup in the config to see if
// there's an override set we should apply
// the hostname used is dervied from [in order]
//    1) the hostnameOverride parameter if not ""
//    2) the value of the Environemnt variable in envKeyName, if not ""
//    3) the OS supplied hostname
// if the supplied hostname doesn't exist and is not a fully qualified name
// and there's an entry in hosts that is fully qualified, that'll be returned
func (configs *Configurations) resolveHostname(envKeyName, hostnameOverride string) (string, error) {
	var err error
	hn := hostnameOverride
	if hn == "" {
		if envKeyName != "" {
			hn = os.Getenv(envKeyName)
		}
		if hn == "" {
			if hn, err = os.Hostname(); err != nil {
				return "", err
			}
		}
	}
	if _, exists := configs.Hosts[hn]; !exists {
		// resolved host name doesn't appear in the Hosts section, see if
		// host is not fully qualified, and see if there's a FQ version
		// in the host list.
		if strings.Index(hn, ".") == -1 {
			// no quick way to do this, other than to trawl through them all
			qualhn := hn + "."
			for k := range configs.Hosts {
				if strings.HasPrefix(k, qualhn) {
					return k, nil
				}
			}
		}
	}
	return hn, nil
}
