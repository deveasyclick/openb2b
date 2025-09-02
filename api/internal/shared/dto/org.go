package dto

import (
	"github.com/deveasyclick/openb2b/internal/model"
)

// CreateOrgDTO represents the payload for creating a new organization
// @Description Organization creation request
type CreateOrgDTO struct {
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

	Address AddressRequired `json:"address"`
}

// UpdateOrgDTO represents the payload for updating an organization
// @Description Organization update request
type UpdateOrgDTO struct {
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
	Address AddressOptional `json:"address"`
}

func (dto *CreateOrgDTO) ToModel() *model.Org {
	return &model.Org{
		Name:             dto.Name,
		Logo:             dto.Logo,
		OrganizationName: dto.OrganizationName,
		OrganizationUrl:  dto.OrganizationUrl,
		Email:            dto.Email,
		Phone:            dto.Phone,
		Address:          dto.Address.ToModel(),
	}
}

func (dto *UpdateOrgDTO) ApplyModel(org *model.Org) {
	if dto.Name != "" {
		org.Name = dto.Name
	}
	if dto.Logo != "" {
		org.Logo = dto.Logo
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
	dto.Address.ApplyModel(org.Address)
}
