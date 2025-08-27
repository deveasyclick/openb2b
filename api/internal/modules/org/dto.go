package org

import "github.com/deveasyclick/openb2b/internal/model"

// Address represents the address of the organization
// @Description Address
type Address struct {
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

// createDTO represents the payload for creating a new organization
// @Description Organization creation request
type createDTO struct {
	// Name of the organization (short name / brand)
	// Required: true
	// Min length: 3
	// Max length: 50
	Name string `json:"name" validate:"required,min=3,max=50" example:"OpenB2B"`

	// Logo URL of the organization
	// Required: false
	Logo string `json:"logo" validate:"omitempty,url" example:"https://example.com/logo.png"`

	// Full legal organization name
	// Required: true
	OrganizationName string `json:"organizationName" validate:"required,min=3,max=50" example:"OpenB2B Technologies Inc."`

	// Official website URL
	// Required: false
	OrganizationUrl string `json:"organizationUrl" validate:"omitempty,url" example:"https://openb2b.io"`

	// Contact email
	// Required: true
	Email string `json:"email" validate:"required,email" example:"contact@openb2b.io"`

	// Contact phone
	// Required: true
	Phone string `json:"phone" validate:"required,min=10,max=50" example:"+1-202-555-0199"`

	Address Address `json:"address"`
}

// updateDTO represents the payload for updating an organization
// @Description Organization update request
type updateDTO struct {
	// Name of the organization
	Name string `json:"name" validate:"omitempty,min=3,max=50" example:"OpenB2B"`

	// Logo URL of the organization
	Logo string `json:"logo" validate:"omitempty,url" example:"https://example.com/logo.png"`

	// Full legal organization name
	OrganizationName string `json:"organizationName" validate:"omitempty,min=3,max=50" example:"OpenB2B Technologies Inc."`

	// Official website URL
	OrganizationUrl string `json:"organizationUrl" validate:"omitempty,url" example:"https://openb2b.io"`

	// Contact email
	Email string `json:"email" validate:"omitempty,email" example:"contact@openb2b.io"`

	// Contact phone
	Phone string `json:"phone" validate:"omitempty,min=10,max=50" example:"+1-202-555-0199"`

	// Address of the organization
	Address Address `json:"address" validate:"optional"`
}

func (dto *createDTO) ToModel() *model.Org {
	return &model.Org{
		Name:             dto.Name,
		Logo:             dto.Logo,
		OrganizationName: dto.OrganizationName,
		OrganizationUrl:  dto.OrganizationUrl,
		Email:            dto.Email,
		Phone:            dto.Phone,
		Address: &model.Address{
			City:    dto.Address.City,
			State:   dto.Address.State,
			Zip:     dto.Address.Zip,
			Country: dto.Address.Country,
			Address: dto.Address.Address,
		},
	}
}

func (dto *updateDTO) ApplyModel(org *model.Org) {
	if dto.Name != "" {
		org.Name = dto.Name
	}
	if dto.Logo != "" {
		org.Logo = dto.Logo
	}
	if dto.OrganizationName != "" {
		org.OrganizationName = dto.OrganizationName
	}
	if dto.OrganizationUrl != "" {
		org.OrganizationUrl = dto.OrganizationUrl
	}
	if dto.Email != "" {
		org.Email = dto.Email
	}
	if dto.Phone != "" {
		org.Phone = dto.Phone
	}
	if dto.Address.State != "" {
		org.Address.State = dto.Address.State
	}
	if dto.Address.Address != "" {
		org.Address.Address = dto.Address.Address
	}
	if dto.Address.City != "" {
		org.Address.City = dto.Address.City
	}
	if dto.Address.Country != "" {
		org.Address.Country = dto.Address.Country
	}
	if dto.Address.Zip != "" {
		org.Address.Country = dto.Address.Zip
	}
}
