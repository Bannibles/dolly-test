package v1

// CertCentralBindingOrg provides Organization info for the binded key
type CertCentralBindingOrg struct {
	// ID of the organization
	ID int `json:"organization_id"`
}

// CertCentralBindingProduct provides available Product info for the binded key
type CertCentralBindingProduct struct {
	// NameID of the Certificate order type
	NameID string `json:"name_id"`
}

// CertCentralAccountBinding provides binding between CertCentral and ACME Proxy
type CertCentralAccountBinding struct {
	Name         string                    `json:"name,omitempty"`
	Key          string                    `json:"key"`
	Organization CertCentralBindingOrg     `json:"organization"`
	Product      CertCentralBindingProduct `json:"product"`
}

// CertCentralRegisterKeyRequest specifies request to register CertCentral API key
type CertCentralRegisterKeyRequest CertCentralAccountBinding

// CertCentralRegisterKeyResponse provides response for CertCentralRegisterAPIKeyRequest
type CertCentralRegisterKeyResponse struct {
	// AccountID specifies registered identifier for the API key
	AccountID string `json:"account_id"`
	// DirectoryURL is ACME's URL for the registered API Key
	DirectoryURL string `json:"directory"`
}

// CertCentralGetKeyInfoRequest specifies request to return CertCentral account info,
// that was created by CertCentralRegisterKeyResponse
type CertCentralGetKeyInfoRequest struct{}

// CertCentralGetKeyInfoResponse provides response for CertCentralGetKeyInfoRequest
type CertCentralGetKeyInfoResponse struct {
	// AccountID specifies registered identifier for the API key
	AccountID    string                    `json:"account_id"`
	Name         string                    `json:"name,omitempty"`
	Organization CertCentralBindingOrg     `json:"organization"`
	Product      CertCentralBindingProduct `json:"product"`
}

// CertCentralUnregisterKeyRequest specifies request to unregister CertCentral account,
// that was created by CertCentralRegisterKeyResponse
type CertCentralUnregisterKeyRequest struct{}

// CertCentralUnregisterKeyResponse provides response for CertCentralUnregisterKeyRequest
type CertCentralUnregisterKeyResponse struct {
	// AccountID specifies registered identifier for the API key
	AccountID    string                    `json:"account_id"`
	Name         string                    `json:"name,omitempty"`
	Organization CertCentralBindingOrg     `json:"organization"`
	Product      CertCentralBindingProduct `json:"product"`
	// Deleted specifies number of records deleted
	Deleted int64 `json:"deleted"`
}
