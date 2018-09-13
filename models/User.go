package models

// User structure
// The firebase user also contains 'emailVerified' and 'photoURL'. The first is default false and should not
// be possible to set through a request. Our platform has no user images so ignored for now.
type User struct {
	Uid string `json:"uid"`
	Email string `json:"email"`
	EmailVerified bool `json:"emailVerified"`
	PhoneNumber string `json:"phoneNumber"`
	Password string `json:"password"`
	DisplayName string `json:"displayName"`
	Disabled bool `json:"disabled"`
}