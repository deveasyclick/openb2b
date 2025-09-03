package dto

import "github.com/deveasyclick/openb2b/internal/model"

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

func (a *AddressOptional) ToModel() *model.Address {
	return &model.Address{
		Address: a.Address,
		City:    a.City,
		State:   a.State,
		Country: a.Country,
		Zip:     a.Zip,
	}
}

func (a *AddressOptional) ApplyModel(address *model.Address) {
	if a.Address != "" {
		address.Address = a.Address
	}
	if a.City != "" {
		address.City = a.City
	}
	if a.State != "" {
		address.State = a.State
	}
	if a.Country != "" {
		address.Country = a.Country
	}
	if a.Zip != "" {
		address.Zip = a.Zip
	}
}

func (a *AddressRequired) ToModel() *model.Address {
	return &model.Address{
		Address: a.Address,
		City:    a.City,
		State:   a.State,
		Country: a.Country,
		Zip:     a.Zip,
	}
}

func (a *AddressRequired) ApplyModel(address *model.Address) {
	if a.Address != "" {
		address.Address = a.Address
	}
	if a.City != "" {
		address.City = a.City
	}
	if a.State != "" {
		address.State = a.State
	}
	if a.Country != "" {
		address.Country = a.Country
	}
	if a.Zip != "" {
		address.Zip = a.Zip
	}
}
