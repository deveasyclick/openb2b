package types

// Email represents an email address structure
type ClerkEmail struct {
	ID           string `json:"id"`
	EmailAddress string `json:"email_address"`
	Verified     bool   `json:"verified"`
}

// ClerkUser represents the Clerk user data structure
type ClerkUser struct {
	ID             string       `json:"id"`
	FirstName      string       `json:"first_name"`
	LastName       string       `json:"last_name"`
	EmailAddresses []ClerkEmail `json:"email_addresses"`
}
