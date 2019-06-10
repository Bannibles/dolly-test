package jwtmapper_test

import (
	"net/http"
	"testing"
	"time"

	"github.com/certcentral/enrollme/api/v1"
	"github.com/certcentral/enrollme/pkg/roles/jwtmapper"
	"github.com/go-phorce/dolly/xhttp/header"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_Config(t *testing.T) {
	_, err := jwtmapper.LoadConfig("testdata/missing.yaml")
	require.Error(t, err)
	assert.Equal(t, "open testdata/missing.yaml: no such file or directory", err.Error())

	_, err = jwtmapper.LoadConfig("testdata/roles_corrupted.1.yaml")
	require.Error(t, err)
	assert.Equal(t, `unable to unmarshal "testdata/roles_corrupted.1.yaml": yaml: line 2: mapping values are not allowed in this context`, err.Error())

	_, err = jwtmapper.LoadConfig("testdata/roles_corrupted.2.yaml")
	require.Error(t, err)
	assert.Equal(t, `unable to unmarshal "testdata/roles_corrupted.2.yaml": yaml: line 2: did not find expected key`, err.Error())

	_, err = jwtmapper.LoadConfig("testdata/roles_no_kid.yaml")
	require.Error(t, err)
	assert.Equal(t, `missing kid: "testdata/roles_no_kid.yaml"`, err.Error())

	_, err = jwtmapper.LoadConfig("testdata/roles_no_keys.yaml")
	require.Error(t, err)
	assert.Equal(t, `missing keys: "testdata/roles_no_keys.yaml"`, err.Error())

	cfg, err := jwtmapper.LoadConfig("testdata/roles.yaml")
	require.NoError(t, err)
	require.NotNil(t, cfg)
	assert.Equal(t, 2, len(cfg.Keys))
	assert.Equal(t, "1", cfg.KeyID)
	assert.Equal(t, "enrollme", cfg.Issuer)
	assert.Equal(t, "enrollme", cfg.Audience)
	assert.Equal(t, "enrollme-user", cfg.DefaultRole)
	assert.Equal(t, 2, len(cfg.RolesMap["enrollme-admin"]))
}

func Test_Load(t *testing.T) {
	_, err := jwtmapper.Load("testdata/missing.yaml")
	require.Error(t, err)
	assert.Equal(t, "open testdata/missing.yaml: no such file or directory", err.Error())

	_, err = jwtmapper.Load("testdata/roles_corrupted.1.yaml")
	require.Error(t, err)

	_, err = jwtmapper.Load("testdata/roles_corrupted.2.yaml")
	require.Error(t, err)

	_, err = jwtmapper.Load("testdata/roles.yaml")
	require.NoError(t, err)

	_, err = jwtmapper.Load("")
	require.NoError(t, err)
}

func Test_Sign(t *testing.T) {
	p, err := jwtmapper.Load("testdata/roles.yaml")
	require.NoError(t, err)
	p1, err := jwtmapper.Load("testdata/roles.1.yaml")
	require.NoError(t, err)
	p2, err := jwtmapper.Load("testdata/roles.2.yaml")
	require.NoError(t, err)

	t.Run("default role", func(t *testing.T) {
		userInfo := &v1.UserInfo{
			ID:    "123",
			OrgID: "456",
			Email: "denis@ekspand.com",
		}

		auth, err := p.SignToken(userInfo, "device123", time.Minute)
		require.NoError(t, err)

		r, _ := http.NewRequest(http.MethodGet, "/", nil)
		setAuthorizationHeader(r, auth.AccessToken, "device123")
		assert.True(t, p.Applicable(r))

		id, err := p.IdentityMapper(r)
		require.NoError(t, err)
		assert.Equal(t, "enrollme-admin", id.Role())
		assert.Equal(t, userInfo.Email, id.Name())
	})

	t.Run("admin role", func(t *testing.T) {
		userInfo := &v1.UserInfo{
			ID:    "123",
			OrgID: "456",
			Email: "denis@ekspand.com",
		}

		auth, err := p.SignToken(userInfo, "device123", time.Minute)
		require.NoError(t, err)

		r, _ := http.NewRequest(http.MethodGet, "/", nil)
		setAuthorizationHeader(r, auth.AccessToken, "device123")
		assert.True(t, p.Applicable(r))

		id, err := p.IdentityMapper(r)
		require.NoError(t, err)
		assert.Equal(t, "enrollme-admin", id.Role())
		assert.Equal(t, userInfo.Email, id.Name())
	})

	t.Run("invalid_sig", func(t *testing.T) {
		userInfo := &v1.UserInfo{
			ID:    "123",
			OrgID: "456",
			Email: "denis@ekspand.com",
		}

		auth, err := p.SignToken(userInfo, "device123", time.Minute)
		require.NoError(t, err)

		r, _ := http.NewRequest(http.MethodGet, "/", nil)
		setAuthorizationHeader(r, auth.AccessToken, "device123")

		assert.True(t, p2.Applicable(r))
		_, err = p2.IdentityMapper(r)
		require.Error(t, err)
		assert.Equal(t, "failed to verify token: signature is invalid", err.Error())
	})

	t.Run("invalid_issuer", func(t *testing.T) {
		userInfo := &v1.UserInfo{
			ID:    "123",
			OrgID: "456",
			Email: "denis@ekspand.com",
		}

		auth, err := p.SignToken(userInfo, "device123", time.Minute)
		require.NoError(t, err)

		r, _ := http.NewRequest(http.MethodGet, "/", nil)
		setAuthorizationHeader(r, auth.AccessToken, "device123")

		assert.True(t, p1.Applicable(r))
		_, err = p1.IdentityMapper(r)
		require.Error(t, err)
		assert.Equal(t, "invalid issuer: enrollme", err.Error())
	})
}

func Test_Verify(t *testing.T) {
	p, err := jwtmapper.Load("testdata/roles.yaml")
	require.NoError(t, err)

	userInfo := &v1.UserInfo{
		ID:    "123",
		OrgID: "456",
		Email: "denis@ekspand.com",
	}
	auth, err := p.SignToken(userInfo, "device123", time.Second)
	require.NoError(t, err)

	t.Run("invalid_token", func(t *testing.T) {
		r, _ := http.NewRequest(http.MethodGet, "/", nil)
		setAuthorizationHeader(r, auth.AccessToken+"123", "device123")
		_, err = p.IdentityMapper(r)
		require.Error(t, err)
		assert.Equal(t, "failed to verify token: signature is invalid", err.Error())
	})

	t.Run("invalid_deviceID", func(t *testing.T) {
		r, _ := http.NewRequest(http.MethodGet, "/", nil)
		setAuthorizationHeader(r, auth.AccessToken, "device456")
		_, err = p.IdentityMapper(r)
		require.Error(t, err)
		assert.Equal(t, "invalid deviceID: device456", err.Error())
	})

	t.Run("not_applicable_header", func(t *testing.T) {
		r, _ := http.NewRequest(http.MethodGet, "/", nil)
		r.Header.Set("Authorization", auth.AccessToken)
		id, err := p.IdentityMapper(r)
		require.NoError(t, err)
		assert.Nil(t, id)
	})
}

// setAuthorizationHeader applies Authorization header
func setAuthorizationHeader(r *http.Request, token, deviceID string) {
	r.Header.Set(header.Authorization, header.Bearer+" "+token)
	if deviceID != "" {
		r.Header.Set(header.XDeviceID, deviceID)
	}
}
