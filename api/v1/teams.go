package v1

import (
	"time"
)

const (
	// MaxUserNameLen specifies maximum length for user's name
	MaxUserNameLen = 64
	// MaxEmailNameLen specifies maximum length for email
	MaxEmailNameLen = 64
	// MaxTeamNameLen specifies maximum length for team's name
	MaxTeamNameLen = 64
)

// User provides basic user information
type User struct {
	ID          string     `json:"id"`
	Name        string     `json:"name"`
	Email       string     `json:"email"`
	Age         int        `json:"age"`
	LoginCount  int        `json:"login_count"`
	LastLoginAt *time.Time `json:"last_login_at,omitempty"`
}

// Team provides basic team information
type Team struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
}

// TeamMembership provides team membership information for a user
type TeamMembership struct {
	ID     string `json:"id"`
	TeamID string `json:"team_id"`
	Team   string `json:"team"`
	UserID string `json:"user_id"`
	Role   string `json:"role"`
	Age    int    `json:"age"`
}

// TeamMemberInfo provides team membership information for a user
type TeamMemberInfo struct {
	MembershipID string `json:"membership_id"`
	TeamID       string `json:"team_id"`
	Team         string `json:"team"`
	UserID       string `json:"user_id"`
	Role         string `json:"role"`
	Name         string `json:"name"`
	Email        string `json:"email"`
	Age          int    `json:"age"`
}

// GetTeamMembershipsResponse returns teams membership
type GetTeamMembershipsResponse struct {
	Memberships []*TeamMemberInfo `json:"memberships"`
}

// ListTeamsResponse returns list of available teams
type ListTeamsResponse struct {
	Teams []string `json:"teams"`
}

// FindUserRequest specifies user search request
type FindUserRequest struct {
	Name   string `json:"name"`
	MinAge int    `json:"min_age"`
	MaxAge int    `json:"max_age"`
}

// FindUserResponse returns list of users that match the search criteria
type FindUserResponse struct {
	Users []*User `json:"users"`
}
