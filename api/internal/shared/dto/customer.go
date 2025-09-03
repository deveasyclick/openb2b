package dto

import "github.com/deveasyclick/openb2b/internal/model"

// CreateCustomerDTO represents incoming/outgoing API data for Customer
type CreateCustomerDTO struct {
	FirstName   string           `json:"firstName" validate:"required,max=100"`
	LastName    string           `json:"lastName" validate:"required,max=100"`
	PhoneNumber string           `json:"phoneNumber" validate:"required"`
	Email       string           `json:"email,omitempty"`
	Address     *AddressOptional `json:"address,omitempty"`
	Company     string           `json:"company,omitempty"`
}

// ToModel converts CreateCustomerDTO to a Customer model
func (dto *CreateCustomerDTO) ToModel(orgID uint) *model.Customer {
	customer := &model.Customer{
		FirstName:   dto.FirstName,
		LastName:    dto.LastName,
		PhoneNumber: dto.PhoneNumber,
		Email:       dto.Email,
		Company:     dto.Company,
		OrgID:       orgID,
	}

	if dto.Address != nil {
		customer.Address = dto.Address.ToModel()
	}
	return customer
}

type UpdateCustomerDTO struct {
	FirstName *string          `json:"firstName" validate:"omitempty,max=100"`
	LastName  *string          `json:"lastName" validate:"omitempty,max=100"`
	Email     *string          `json:"email" validate:"omitempty,max=100"`
	Company   *string          `json:"company" validate:"omitempty,max=100"`
	Address   *AddressOptional `json:"address,omitempty"`
}

// ApplyModel updates an existing Customer model with DTO values
func (dto *UpdateCustomerDTO) ApplyModel(c *model.Customer) {
	if dto.FirstName != nil {
		c.FirstName = *dto.FirstName
	}
	if dto.LastName != nil {
		c.LastName = *dto.LastName
	}

	if dto.Company != nil {
		c.Company = *dto.Company
	}

	if dto.Address != nil {
		dto.Address.ApplyModel(c.Address)
	}

	if dto.Email != nil {
		c.Email = *dto.Email
	}
}
