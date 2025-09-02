package dto

// Address represents the address of the organization
// @Description Address
type AddressRequired struct {
	// State where the organization is located
	// Required: true
	State string `json:"state" validate:"required,min=2,max=30" example:"California"`

	// City where the organization is located
	// Required: true
	City string `json:"city" validate:"required,min=2,max=30" example:"San Francisco"`

	// Address of the organization
	// Required: true
	Address string `json:"address" validate:"required,min=5,max=100" example:"123 Market Street"`

	// Country where the organization is registered
	// Required: true
	Country string `json:"country" validate:"required,min=2,max=100" example:"USA"`

	// Zip where the organization is located
	// Required: true
	Zip string `json:"zip" validate:"required,min=2,max=30" example:"02912"`
}

type AddressOptional struct {
	// State where the organization is located
	// Required: false
	State string `json:"state" validate:"omitempty,min=2,max=30" example:"California"`

	// City where the organization is located
	// Required: false
	City string `json:"city" validate:"omitempty,min=2,max=30" example:"San Francisco"`

	// Address of the organization
	// Required: false
	Address string `json:"address" validate:"omitempty,min=5,max=100" example:"123 Market Street"`

	// Country where the organization is registered
	// Required: false
	Country string `json:"country" validate:"omitempty,min=2,max=100" example:"USA"`

	// Zip where the organization is located
	// Required: false
	Zip string `json:"zip" validate:"omitempty,min=2,max=30" example:"02912"`
}
