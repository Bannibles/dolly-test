package v1

// UserInfo provides basic info about user
// This structure is returned by Aloha, so don't change!
type UserInfo struct {
	ID                string `json:"user_id"`
	OrgID             string `json:"organization_id"`
	PreferredUserName string `json:"preferred_username"`
	Name              string `json:"name"`
	Email             string `json:"email"`
	EmailVerified     bool   `json:"email_verified"`
	ProfileURL        string `json:"profile"`
	PictureURL        string `json:"picture"`
	Phone             string `json:"phone_number"`
	PhoneVerified     bool   `json:"phone_number_verified"`
	Active            bool   `json:"active"`
	UserType          string `json:"user_type"`
	TimeZone          string `json:"zoneinfo"`
	UTCOffset         int    `json:"utcOffset"`
}
