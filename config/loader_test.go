package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/juju/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const projFolder = ".."

func Test_trimAll(t *testing.T) {
	test := func(src []string, exp []string) {
		res := trimAll(src)
		assert.Equal(t, exp, res, "trimAll(%v) unexpected result")
	}
	test([]string{":500", ":501"}, []string{":500", ":501"})
	test([]string{"  :500 ", " :501 "}, []string{":500", ":501"})
	test([]string{"  :500 ", "", " :501 "}, []string{":500", ":501"})
	test([]string{"  :500 ", "    ", " ", " "}, []string{":500"})
	test([]string{}, []string{})
	test([]string{"", " ", " ", " ", "   "}, []string{})
}

func Test_ConfigFilesAreJson(t *testing.T) {
	isJSON := func(file string) {
		abs := "../" + file
		f, err := os.Open(abs)
		require.NoError(t, err, "Unable to open file: %v", file)
		defer f.Close()
		var v map[string]interface{}
		assert.NoError(t, json.NewDecoder(f).Decode(&v), "JSON parser error for file %v", file)
	}
	isJSON("etc/dev/" + ConfigFileName)
}

func Test_LoadConfig(t *testing.T) {
	_, _, err := LoadConfig("missing.json")
	assert.Error(t, err)
	assert.True(t, errors.IsNotFound(err) || os.IsNotExist(err), "LoadConfig with missing file should return a file doesn't exist error: %v", errors.Trace(err))

	cfgFile, err := GetConfigAbsFilename("etc/dev/"+ConfigFileName, projFolder)
	require.NoError(t, err, "unable to determine config file")

	c, _, err := LoadConfig(cfgFile)
	require.NoError(t, err, "failed to load config: %v", cfgFile)

	testDirAbs := func(name, dir string) {
		if dir != "" {
			assert.True(t, filepath.IsAbs(dir), "dir %q should be an absoluite path", name)
		}
	}
	testDirAbs("HTTPS.ServerTLS.CertFile", c.HTTPS.ServerTLS.CertFile)
	testDirAbs("HTTPS.ServerTLS.KeyFile", c.HTTPS.ServerTLS.KeyFile)
}
