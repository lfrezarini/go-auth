package jsonwebtoken

// Claims represents the claims that can be passed on the creation of jsonwebtoken
type Claims struct {
	//Iss is the issuer of the JWT (the api that generated it).
	Iss string `json:"iss"`

	//Sub is the subject that generated the token
	Sub string `json:"sub"`
}
