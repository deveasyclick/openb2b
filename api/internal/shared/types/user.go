package types

// Email represents an email address structure
type ClerkEmail struct {
	ID           string `json:"id" mapstructure:"id"`
	EmailAddress string `json:"email_address" mapstructure:"email_address"`
	Verified     bool   `json:"verified" mapstructure:"verified"`
}

// ClerkUser represents the Clerk user data structure
type ClerkUser struct {
	ID             string       `mapstructure:"id" json:"id"`
	FirstName      string       `mapstructure:"first_name" json:"first_name"`
	LastName       string       `mapstructure:"last_name" json:"last_name"`
	EmailAddresses []ClerkEmail `mapstructure:"email_addresses" json:"email_addresses"`
}
