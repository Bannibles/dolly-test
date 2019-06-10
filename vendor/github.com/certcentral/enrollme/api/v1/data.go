package v1

// GetDataKeysRequest specifies a POST request to retrieve DB Keys
type GetDataKeysRequest struct {
	// Limit specifies maximum number of keys to return.
	// If this parameter is not set, then a default value will be used [1000 records]
	Limit int64 `json:"limit"`

	// CountOnly specifies that values shall not be returned.
	// If this parameter is true, then Limit is ignored
	CountOnly bool `json:"count"`

	// Start is used for pagination.
	// If this parameter is set, then the keys retrieved starting from specified key,
	// but not including it
	Start string `json:"start"`

	// End is used to specify the range of the keys in query.
	// If this parameter is empty, then all keys from start will be searched.
	End string `json:"end"`
}

// GetDataKeysResponse specifies a response for GetDataKeysRequest
type GetDataKeysResponse struct {
	Keys []string `json:"keys"`

	// LastKey specifies the name of the last key to be used for subsequent requests
	LastKey string `json:"last"`

	// More specifies if there are more records to retrieve
	More bool `json:"more"`

	// Count specifies total number of records found for the request
	Count int64 `json:"count"`
}

// GetUpdateDataEncryptionKeysResponse specifies a response for /v1/data/update_cek
type GetUpdateDataEncryptionKeysResponse struct {
	UpdatedNodes []string `json:"updated_nodes"`
}
