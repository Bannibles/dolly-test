package v1

// Public API
const (
	// URIForTeams returns teams
	//
	// Verbs: GET
	URIForTeams = "/v1/teams"

	// URIForTeamsMemberships returns teams membership for the caller
	//
	// Verbs: GET
	URIForTeamsMemberships = URIForTeams + "/memberships"

	// URIForUsers returns users
	//
	// Verbs: GET
	// Parameters:
	//	name		- optional, name of the user to filter by
	//  max_age		- optional, max age of the user to filter by
	//  min_age		- optional, min age of the user to filter by
	URIForUsers = "/v1/users"
)
