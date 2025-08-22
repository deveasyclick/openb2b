package org

import "github.com/deveasyclick/openb2b/internal/model"

type CreateOrgDTO struct {
	Name             string `json:"name" validate:"required,min=3,max=50"`
	Logo             string `json:"logo" validate:"omitempty,url"`
	OrganizationName string `json:"organizationName" validate:"required,min=3,max=50"`
	OrganizationUrl  string `json:"organizationUrl" validate:"omitempty,url"`
	Email            string `json:"email" validate:"required,email"`
	Phone            string `json:"phone" validate:"required,min=10,max=50"`
	State            string `json:"state" validate:"required,min=2,max=30"`
	City             string `json:"city" validate:"required,min=2,max=30"`
	Address          string `json:"address" validate:"required,min=5,max=100"`
	Country          string `json:"country" validate:"required,min=2,max=100"`
}

type UpdateOrgDTO struct {
	Name             string `json:"name" validate:"omitempty,min=3,max=50"`
	Logo             string `json:"logo" validate:"omitempty,url"`
	OrganizationName string `json:"organizationName" validate:"omitempty,min=3,max=50"`
	OrganizationUrl  string `json:"organizationUrl" validate:"omitempty,url"`
	Email            string `json:"email" validate:"omitempty,email"`
	Phone            string `json:"phone" validate:"omitempty,min=10,max=50"`
	State            string `json:"state" validate:"omitempty,min=2,max=30"`
	City             string `json:"city" validate:"omitempty,min=2,max=30"`
	Address          string `json:"address" validate:"omitempty,min=5,max=100"`
}

func (dto *CreateOrgDTO) ToModel() *model.Org {
	return &model.Org{
		Name:             dto.Name,
		Logo:             dto.Logo,
		OrganizationName: dto.OrganizationName,
		OrganizationUrl:  dto.OrganizationUrl,
		Email:            dto.Email,
		Phone:            dto.Phone,
		State:            dto.State,
		City:             dto.City,
		Address:          dto.Address,
		Country:          dto.Country,
	}
}

func (dto *UpdateOrgDTO) ApplyModel(org *model.Org) {
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
	if dto.State != "" {
		org.State = dto.State
	}
	if dto.City != "" {
		org.City = dto.City
	}
	if dto.Address != "" {
		org.Address = dto.Address
	}
}
