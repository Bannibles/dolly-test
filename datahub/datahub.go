package datahub

import (
	"context"

	"github.com/go-phorce/dolly-test/api/v1"
)

// UsersManager interface provides sample user management API
type UsersManager interface {
	ListTeams(ctx context.Context) (*v1.ListTeamsResponse, error)
	FindUser(ctx context.Context, req *v1.FindUserRequest) (*v1.FindUserResponse, error)
}

// Datahub defines an interface to work with data storage
type Datahub interface {
	UsersManager
}
