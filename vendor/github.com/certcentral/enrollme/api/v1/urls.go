package v1

// Status service API
const (
	// URIForStatus provides information about the server
	//
	// Verbs: GET
	URIForStatus = "/v1/status"

	// URIForStatusCluster provides information about the cluster nodes
	//
	// Verbs: GET
	URIForStatusCluster = "/v1/status/cluster"

	// URIForStatusServer provides information about current node status
	//
	// Verbs: GET
	URIForStatusServer = "/v1/status/server"

	// URIForStatusCerts returns a requested server certificate
	//
	// Verbs: GET
	//
	// Parameters:
	//	type	- type of the certificate: server|peer|tls
	URIForStatusCerts = "/v1/status/certs/:type"

	// URIForVersion provides information about current version of the service
	//
	// Verbs: GET
	URIForVersion = "/v1/version"
)

// Cluster management API
const (
	URIForCluster = "/v1/cluster"

	// URIForClusterNodes adds or removes a specified member from the cluster
	//
	// Verbs: GET, POST, DELETE
	URIForClusterNodes            = URIForCluster + "/nodes"
	URIForClusterNodeTemplateByID = URIForClusterNodes + "/:node_id"

	// URIForClusterHeartbeat accepts POST requests with heartbeat from a service
	//
	// Verbs: POST
	URIForClusterHeartbeat = URIForCluster + "/heartbeat"

	// URIForClusterNonce provides nonce service for the cluster
	//
	// Verbs: HEAD, GET
	//
	// HEAD will generate nonce and return "Replay-Nonce" header
	// GET will verify nonce and mark it as used.
	URIForClusterNonce             = URIForCluster + "/nonce"
	URIForClusterNonceTemplateByID = URIForClusterNonce + "/:nonce"
)

// DB management API
const (
	URIForDB = "/v1/data"

	// URIForDBListKeys provides pagination API to retrieve keys in DB
	//
	// Verbs: POST
	URIForDBListKeys = "/v1/data/list_keys"

	// URIForDBUpdateCEK tries to update Content Encryption Keys for new cluster nodes
	//
	// Verbs: GET
	URIForDBUpdateCEK = "/v1/data/update_cek"
)

// CertCentral integration API
const (
	URIForCertCentral = "/v1/certcentral"

	// URIForCertCentralBinding supports account binding
	//
	// Verbs:
	//    POST CertCentralAPIKeyInfo
	URIForCertCentralBinding = "/v1/certcentral/binding"

	// URIForCertCentralBindingByID provides binding management for specific account
	// Verbs:
	//    GET
	//    DELETE
	URIForCertCentralBindingByID = "/v1/certcentral/binding/:id"
)

// Auth API
const (
	URIForAuth = "/v1/auth"

	// URIForAuthStsURL provides auth URL
	//
	// Verbs: GET
	// Paramaters:
	//	redirect_url
	//	device_id
	URIForAuthStsURL = URIForAuth + "/sts/url"

	// URIForAuthTokenRefresh provides auth token refresh
	//
	// Verbs: GET
	URIForAuthTokenRefresh = URIForAuth + "/token/refresh"

	// URIForAuthCallback provides auth URL
	//
	// Verbs: GET
	URIForAuthCallback = URIForAuth + "/token/callback"
)
