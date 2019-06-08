package service_test

import (
	"testing"

	"github.com/go-phorce/dolly-test/service/teams"
	"github.com/go-phorce/dolly/rest"
	"github.com/stretchr/testify/require"
)

var serviceFactories = map[string]func(server rest.Server) interface{}{
	teams.ServiceName: teams.Factory,
}

func Test_invalidArgs(t *testing.T) {
	for _, f := range serviceFactories {
		testInvalidServiceArgs(t, f)
	}
}

func testInvalidServiceArgs(t *testing.T, f func(server rest.Server) interface{}) {
	defer func() {
		err := recover()
		require.NotNil(t, err, "Expected panic but didn't get one")
	}()
	f(nil)
}
