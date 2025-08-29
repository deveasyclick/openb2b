package product

import "github.com/deveasyclick/openb2b/internal/model"

type createProductVariantDTO struct {
	SKU   string  `json:"sku" validate:"required,min=2,max=50"`
	Color string  `json:"color" validate:"omitempty,min=1,max=30"`
	Size  string  `json:"size" validate:"omitempty,min=1,max=30"`
	Price float64 `json:"price" validate:"required,gt=0"`
	Stock int     `json:"stock" validate:"required,min=0"`
}

type createProductDTO struct {
	Name        string                    `json:"name" validate:"required,min=2,max=100"`
	Category    string                    `json:"category" validate:"omitempty,min=2,max=50"`
	ImageURL    string                    `json:"imageUrl" validate:"omitempty"`
	Description string                    `json:"description" validate:"omitempty,min=2,max=1000"`
	Variants    []createProductVariantDTO `json:"variants" validate:"required,dive"`
}

func (p *createProductDTO) ToModel(orgID uint) model.Product {
	product := model.Product{
		Name:        p.Name,
		Category:    p.Category,
		ImageURL:    p.ImageURL,
		Description: p.Description,
		OrgID:       orgID,
	}

	// map variants
	for _, v := range p.Variants {
		product.Variants = append(product.Variants, model.Variant{
			SKU:   v.SKU,
			Color: v.Color,
			Size:  v.Size,
			Price: v.Price,
			Stock: v.Stock,
			OrgID: orgID, // enforce same org
		})
	}

	return product
}

// ----- UPDATE -----

type updateProductDTO struct {
	Name        *string `json:"name" validate:"omitempty,min=2,max=100"`
	Category    *string `json:"category" validate:"omitempty,min=2,max=50"`
	ImageURL    *string `json:"imageUrl" validate:"omitempty"`
	Description *string `json:"description" validate:"omitempty,min=2,max=1000"`
}

func (dto *updateProductDTO) ApplyModel(product *model.Product) {
	if dto.Name != nil {
		product.Name = *dto.Name
	}
	if dto.Category != nil {
		product.Category = *dto.Category
	}
	if dto.ImageURL != nil {
		product.ImageURL = *dto.ImageURL
	}
	if dto.Description != nil {
		product.Description = *dto.Description
	}
}

// Variants
type createVariantDTO struct {
	SKU   string  `json:"sku" validate:"required,min=2,max=50"`
	Color string  `json:"color" validate:"omitempty,min=1,max=30"`
	Size  string  `json:"size" validate:"omitempty,min=1,max=30"`
	Price float64 `json:"price" validate:"required,gt=0"`
	Stock int     `json:"stock" validate:"required,min=0"`
}

func (v *createVariantDTO) ToModel() model.Variant {
	return model.Variant{
		SKU:   v.SKU,
		Color: v.Color,
		Size:  v.Size,
		Price: v.Price,
		Stock: v.Stock,
	}
}

type updateVariantDTO struct {
	Color *string  `json:"color" validate:"omitempty,min=1,max=30"`
	Size  *string  `json:"size" validate:"omitempty,min=1,max=30"`
	Price *float64 `json:"price" validate:"omitempty,gt=0"`
	Stock *int     `json:"stock" validate:"omitempty,min=0"`
}

func (dto *updateVariantDTO) ApplyModel(variant *model.Variant) {
	if dto.Color != nil {
		variant.Color = *dto.Color
	}
	if dto.Size != nil {
		variant.Size = *dto.Size
	}
	if dto.Price != nil {
		variant.Price = *dto.Price
	}
	if dto.Stock != nil {
		variant.Stock = *dto.Stock
	}
}
