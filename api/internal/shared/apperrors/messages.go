package apperrors

const (
	// Generic
	ErrEncodeResponse     = "failed to encode response"
	ErrInvalidId          = "invalid id"
	ErrInvalidRequestBody = "invalid request body"
	ErrInvalidFilter      = "invalid filter"

	// Customer
	ErrCustomerNotFound = "customer not found"
	ErrUpdateCustomer   = "error updating customer"
	ErrDeleteCustomer   = "error deleting customer"
	ErrFindCustomer     = "error finding customer"

	// Org
	ErrOrgNotFound      = "org not found"
	ErrUpdateOrg        = "error updating org"
	ErrDeleteOrg        = "error deleting org"
	ErrFindOrg          = "error finding org"
	ErrCreateOrg        = "error creating org"
	ErrOrgAlreadyExists = "org already exists"

	// User
	ErrUserNotFound      = "user not found"
	ErrUpdateUser        = "error updating user"
	ErrDeleteUser        = "error deleting user"
	ErrFindUser          = "error finding user"
	ErrCreateUser        = "error creating user"
	ErrAssignUser        = "error assigning user to org"
	ErrUserAlreadyExists = "user already exists"
	ErrUserFromContext   = "error getting user from context"

	// Product
	ErrProductAlreadyExists = "product already exists"
	ErrCreateProduct        = "error creating product"
	ErrUpdateProduct        = "error updating product"
	ErrDeleteProduct        = "error deleting product"
	ErrFindProduct          = "error finding product"
	ErrProductNotFound      = "product not found"
	ErrFilterProduct        = "error filtering products"

	// Variant
	ErrVariantAlreadyExists = "variant already exists"
	ErrCreateVariant        = "error creating variant"
	ErrUpdateVariant        = "error updating variant"
	ErrDeleteVariant        = "error deleting variant"
	ErrFindVariant          = "error finding variant"
	ErrVariantNotFound      = "variant not found"
)
