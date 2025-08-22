package apperrors

const (
	ErrEncodeResponse = "failed to encode response"

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
	ErrInvalidId   = "invalid id"
)
