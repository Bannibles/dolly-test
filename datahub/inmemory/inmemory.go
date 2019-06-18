package inmemory

import (
	"context"

	"github.com/go-phorce/dolly-test/api/v1"
	"github.com/go-phorce/dolly-test/datahub"
)

type inmem struct {
	teams []string
	users []v1.User
}

// NewUsersManager returns in-memory UsersManager
func NewUsersManager() (datahub.UsersManager, error) {
	p := &inmem{
		teams: []string{"admins", "users"},
		users: []v1.User{
			{ID: "a001", Name: "denis", Email: "denis@ekspand.com", Age: 33},
			{ID: "a002", Name: "andrew", Email: "andrew@ekspand.com", Age: 43},
			{ID: "a003", Name: "hayk", Email: "hayk@ekspand.com", Age: 27},
			{ID: "a004", Name: "daniel", Email: "daniel@ekspand.com", Age: 14},
		},
	}
	return p, nil
}

func (p *inmem) ListTeams(ctx context.Context) (*v1.ListTeamsResponse, error) {
	res := &v1.ListTeamsResponse{
		Teams: p.teams,
	}
	return res, nil
}

func (p *inmem) FindUser(ctx context.Context, req *v1.FindUserRequest) (*v1.FindUserResponse, error) {
	users := make([]*v1.User, 0, len(p.users))

	for idx, u := range p.users {
		if req.Name != "" && u.Name != req.Name {
			// name does not match
			continue
		}

		if req.MinAge > 0 && u.Age < req.MinAge {
			// age does not match
			continue
		}

		if req.MaxAge > 0 && u.Age > req.MaxAge {
			// age does not match
			continue
		}

		users = append(users, &p.users[idx])
	}

	res := &v1.FindUserResponse{
		Users: users,
	}

	return res, nil
}
