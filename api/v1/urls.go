package v1

// Public API
const (
	URIForTeams = "/v1/teams"

	// URIForTeamsMemberships returns teams membership for the caller
	//
	// Verbs: GET
	URIForTeamsMemberships = URIForTeams + "/memberships"
)
