package models

// RefreshToken represents the data structure of a refresh token
type RefreshToken struct {
	// Token is the JWT value
	Token string `json:"token" bson:"token,omitempty"`

	// Identifier is an string used to differentiate the refresh token from others that may exists
	// Ex: the device name, the browser that is using it, etc.
	Identifier string `json:"identifier" bson:"identifier,omitempty"`
}
