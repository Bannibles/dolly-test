package v1

import "time"

// AuthStsURLRequest specifies request to obtain STS URL
type AuthStsURLRequest struct {
	RedirectURL string `json:"redirect_url"`
	DeviceID    string `json:"device_id"`
}

// AuthStsURLResponse provides response for AuthStsURLRequest
type AuthStsURLResponse struct {
	URL string `json:"url"`
}

// UserInfo provides basic info about user
type UserInfo struct {
	ID            string `json:"user_id"`
	OrgID         string `json:"organization_id"`
	Name          string `json:"name"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	ProfileURL    string `json:"profile"`
	PictureURL    string `json:"picture"`
	Phone         string `json:"phone_number"`
	PhoneVerified bool   `json:"phone_number_verified"`
	UserType      string `json:"user_type"`
	TimeZone      string `json:"zoneinfo"`
	UTCOffset     int    `json:"utcOffset"`
}

// Authorization is returned to the client in token refresh response
type Authorization struct {
	Version     string    `json:"version"`
	DeviceID    string    `json:"device_id"`
	UserID      string    `json:"user_id"`
	Email       string    `json:"email"`
	UserName    string    `json:"user_name"`
	Role        string    `json:"role"`
	TokenType   string    `json:"token_type"`
	AccessToken string    `json:"access_token"`
	ExpiresAt   time.Time `json:"expires_at"`
}

// AuthTokenRefreshResponse provides response for token refresh request
type AuthTokenRefreshResponse struct {
	Authorization *Authorization `json:"authorization"`
	Profile       *UserInfo      `json:"profile"`
}
