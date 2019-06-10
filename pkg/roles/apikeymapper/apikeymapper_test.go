package apikeymapper_test

import (
	"net/http"
	"testing"

	"github.com/go-phorce/dolly-test/pkg/roles/apikeymapper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_Config(t *testing.T) {
	_, err := apikeymapper.LoadConfig("testdata/missing.yaml")
	require.Error(t, err)
	assert.Equal(t, "open testdata/missing.yaml: no such file or directory", err.Error())

	_, err = apikeymapper.LoadConfig("testdata/roles_corrupted.1.yaml")
	require.Error(t, err)
	assert.Equal(t, `unable to unmarshal "testdata/roles_corrupted.1.yaml": yaml: line 2: mapping values are not allowed in this context`, err.Error())

	cfg, err := apikeymapper.LoadConfig("testdata/roles_corrupted.2.yaml")
	require.NoError(t, err)
	assert.NotNil(t, 0, len(cfg.KeysMap))

	_, err = apikeymapper.LoadConfig("")
	require.NoError(t, err)

	cfg, err = apikeymapper.LoadConfig("testdata/roles.yaml")
	require.NoError(t, err)
	assert.Equal(t, apikeymapper.APIKeyHeader, cfg.HeaderName)

	cfg, err = apikeymapper.LoadConfig("testdata/roles_default.yaml")
	require.NoError(t, err)
	assert.Equal(t, apikeymapper.APIKeyHeader, cfg.HeaderName)

	p := apikeymapper.New(&apikeymapper.Config{})
	assert.Equal(t, apikeymapper.APIKeyHeader, p.HTTPHeaderName())
}

func Test_identity(t *testing.T) {
	tt := []struct {
		key        string
		identity   string
		experr     string
		applicable bool
	}{
		{key: "549006C71D99A8B14A9DBEEE8319AB6A485F82C922F870F2C3FC7EFDEEB45DE4", identity: "guest/Fred De Gause", applicable: true, experr: ""},
		{key: "A5D461AFA800FBCD820F36707FE9B3051832DB59FD491E6DA91F5E7D884C31C8", identity: "admin/Denis Issoupov", applicable: true, experr: ""},
		{key: "E0D6865224D7A10C533EDB6042FF5603604667952160CD65FDDC170D8F535579", identity: "client/Authenticated User", applicable: true, experr: ""},
		{key: "E0D6865224D7A10C533EDB6042", applicable: true, experr: "invalid access key"},
		{key: "", applicable: false, experr: ""},
	}

	config, err := apikeymapper.LoadConfig("testdata/roles.yaml")
	require.NoError(t, err)
	require.NotNil(t, config)
	require.NotNil(t, config.KeysMap)
	assert.NotNil(t, 3, len(config.KeysMap))

	p := apikeymapper.New(config)
	for _, tc := range tt {
		t.Run(tc.key, func(t *testing.T) {
			r, _ := http.NewRequest(http.MethodGet, "/", nil)
			r.Header.Set(p.HTTPHeaderName(), tc.key)

			applicable := p.Applicable(r)
			assert.Equal(t, tc.applicable, applicable)
			id, err := p.IdentityMapper(r)
			if tc.experr == "" {
				require.NoError(t, err)
				if applicable {
					assert.Equal(t, tc.identity, id.String())
				} else {
					assert.Nil(t, id)
				}
			} else {
				require.Error(t, err)
				assert.Equal(t, tc.experr, err.Error())
			}
		})
	}
}

func Test_Load(t *testing.T) {
	_, err := apikeymapper.Load("testdata/missing.yaml")
	require.Error(t, err)
	assert.Equal(t, "open testdata/missing.yaml: no such file or directory", err.Error())

	_, err = apikeymapper.Load("testdata/roles_corrupted.1.yaml")
	require.Error(t, err)

	_, err = apikeymapper.Load("testdata/roles_corrupted.2.yaml")
	require.NoError(t, err)

	_, err = apikeymapper.Load("testdata/roles.yaml")
	require.NoError(t, err)

	_, err = apikeymapper.Load("")
	require.NoError(t, err)
}
