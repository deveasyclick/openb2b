package apperrors

const (
	// Generic
	ErrEncodeResponse = "failed to encode response"
	ErrInvalidId      = "invalid id"

	// Customer
	ErrCustomerNotFound = "customer not found"
	ErrUpdateCustomer   = "error updating customer"
	ErrDeleteCustomer   = "error deleting customer"
	ErrFindCustomer     = "error finding customer"

	// Org
	ErrOrgNotFound = "org not found"
	ErrUpdateOrg   = "error updating org"
	ErrDeleteOrg   = "error deleting org"
	ErrFindOrg     = "error finding org"
	ErrCreateOrg   = "error creating org"

	// User
	ErrUserNotFound = "user not found"
	ErrUpdateUser   = "error updating user"
	ErrDeleteUser   = "error deleting user"
	ErrFindUser     = "error finding user"
	ErrCreateUser   = "error creating user"
	ErrAssignUser   = "error assigning user to org"
)
