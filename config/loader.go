package config

// Package config allows for the server configuration to read from a separate config file.
// It supports having different configurations for different instance based on host name.
//
// The implementation is primarily provided by the go-phorce/configen tool.

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/go-phorce/dolly/fileutil/resolve"
	"github.com/go-phorce/dolly/netutil"
	"github.com/juju/errors"
)

//go:generate configen -c config_def.json -d . -pkg $TOOLS_PATH

const (
	// ConfigFileName is default name for the configuration file
	ConfigFileName = "demo-config.json"
	envHostnameKey = "DOLLY_HOSTNAME"
)

// Factory is used to create Configuration instance
type Factory struct {
	nodeInfo   netutil.NodeInfo
	searchDirs []string
}

// DefaultFactory returns default configuration factory
func DefaultFactory() (*Factory, error) {
	var nodeInfo netutil.NodeInfo
	nodeInfo, err := netutil.NewNodeInfo(nil)
	if err != nil {
		return nil, errors.Trace(err)
	}

	cwd, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	// try the list of allowed locations to find the config file
	searchDirs := []string{
		filepath.Dir(cwd) + "/etc/dev",
		"/opt/dolly/etc/prod",
		"/var/dolly/etc/prod",
	}

	return &Factory{
		searchDirs: searchDirs,
		nodeInfo:   nodeInfo,
	}, nil
}

// NewFactory returns new configuration factory
func NewFactory(nodeInfo netutil.NodeInfo, searchDirs []string) (*Factory, error) {
	return &Factory{
		searchDirs: searchDirs,
		nodeInfo:   nodeInfo,
	}, nil
}

// LoadConfig will load the server configuration from the named config file,
// apply any overrides, and resolve relative directory locations.
func LoadConfig(configFile string) (*Configuration, string, error) {
	f, err := DefaultFactory()
	if err != nil {
		return nil, "", errors.Trace(err)
	}
	return f.LoadConfigForHostName(configFile, "")
}

// LoadConfig will load the server configuration from the named config file,
// apply any overrides, and resolve relative directory locations.
func (f *Factory) LoadConfig(configFile string) (*Configuration, string, error) {
	return f.LoadConfigForHostName(configFile, "")
}

// LoadConfigForHostName will load the server configuration from the named config file for specified host name,
// apply any overrides, and resolve relative directory locations.
func (f *Factory) LoadConfigForHostName(configFile, hostnameOverride string) (*Configuration, string, error) {
	configFile, baseDir, err := f.resolveConfigFile(configFile)
	if err != nil {
		return nil, "", errors.Trace(err)
	}

	c, err := Load(configFile, envHostnameKey, hostnameOverride)
	if err != nil {
		return nil, "", errors.Trace(err)
	}

	//
	// Substitude ENVIRONMENT
	// Add to this list all configs that require ${NODENAME}, ${HOSTNAME} or ${LOCALIP} substitution
	//
	envVarsResove := []*string{
		&c.HTTPS.ServerTLS.CertFile,
		&c.HTTPS.ServerTLS.KeyFile,
		&c.HTTPS.ServerTLS.TrustedCAFile,
	}
	for _, ptr := range envVarsResove {
		*ptr = f.substitudeEnvVars(*ptr)
	}

	//
	// Resolve Folders and Files
	//

	// Add to this list all configs that require folder resolution to absolute path
	dirsToResove := []*string{
		&c.Audit.Directory,
	}

	filesToResove := []*string{
		&c.HTTPS.ServerTLS.CertFile,
		&c.HTTPS.ServerTLS.KeyFile,
		&c.HTTPS.ServerTLS.TrustedCAFile,
		&c.CryptoProv.Default,
		&c.Authz.CertMapper,
		&c.Authz.APIKeyMapper,
		&c.Authz.JWTMapper,
	}

	optionalFilesToResove := []*string{
		&c.RootCA,
	}

	for _, ptr := range dirsToResove {
		*ptr, err = resolve.Directory(*ptr, baseDir, true)
		if err != nil {
			return nil, "", errors.Annotatef(err, "unable to resolve folder: %s", *ptr)
		}
	}

	for _, ptr := range filesToResove {
		*ptr, err = resolve.File(*ptr, baseDir)
		if err != nil {
			return nil, "", errors.Annotatef(err, "unable to resolve file: %s", *ptr)
		}
	}

	for _, ptr := range optionalFilesToResove {
		*ptr, _ = resolve.File(*ptr, baseDir)
	}

	c.Datacenter = strings.ToLower(c.Datacenter)
	return c, configFile, err
}

// substitudeEnvVars replace ${HOSTNAME}, ${NODENAME} and ${LOCALIP} in input string
func (f *Factory) substitudeEnvVars(s string) string {
	v := strings.Replace(s, "${HOSTNAME}", f.nodeInfo.HostName(), -1)
	v = strings.Replace(v, "${NODENAME}", f.nodeInfo.NodeName(), -1)
	v = strings.Replace(v, "${LOCALIP}", f.nodeInfo.LocalIP(), -1)
	return v
}

// substitudeEnvVarsAll replace ${HOSTNAME}, ${NODENAME} and ${LOCALIP} in input strings
func (f *Factory) substitudeEnvVarsAll(list []string) []string {
	dst := make([]string, len(list))
	for i, s := range list {
		dst[i] = f.substitudeEnvVars(s)
	}
	return dst
}

func (f *Factory) resolveConfigFile(configFile string) (absConfigFile, baseDir string, err error) {
	if configFile == "" {
		configFile = ConfigFileName
	}

	if filepath.IsAbs(configFile) {
		// for absolute, use the folder containing the config file
		baseDir = filepath.Dir(configFile)
		absConfigFile = configFile
		return
	}

	for _, absDir := range f.searchDirs {
		if _, err := os.Stat(absDir); err != nil {
			// folder does not exist or access denied
			continue
		}

		absConfigFile, err = resolve.File(configFile, absDir)
		if err == nil && absConfigFile != "" {
			baseDir = absDir
			return
		}
	}

	err = errors.NotFoundf("file %q in [%s]", configFile, strings.Join(f.searchDirs, ","))
	return
}

func trimAll(p []string) []string {
	emptyCount := 0
	for i := range p {
		p[i] = strings.TrimSpace(p[i])
		if len(p[i]) == 0 {
			emptyCount++
		}
	}
	if emptyCount == 0 {
		return p
	}
	res := make([]string, 0, len(p)-emptyCount)
	for _, q := range p {
		if len(q) > 0 {
			res = append(res, q)
		}
	}
	return res
}

// GetConfigAbsFilename returns absolute path for the configuration file
// from the relative path to projFolder
func GetConfigAbsFilename(file, projFolder string) (string, error) {
	wd, err := os.Getwd() // package dir
	if err != nil {
		return "", errors.Annotate(err, "unable to determine current directory")
	}

	etcCfg, err := filepath.Abs(filepath.Join(wd, projFolder))
	if err != nil {
		return "", errors.Annotatef(err, "unable to determine project directory: %q", projFolder)
	}

	return filepath.Join(etcCfg, file), nil
}
