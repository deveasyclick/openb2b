package model

// Product represents an product entity
// @Description Product response model
type Product struct {
	BaseModel
	Name        string    `gorm:"not null" json:"name"`
	Category    string    `json:"category"`
	OrgID       uint      `gorm:"index;not null" json:"orgId"`
	Org         *Org      `gorm:"foreignKey:OrgID" json:"org,omitempty"`
	ImageURL    string    `json:"imageUrl"`
	Description string    `json:"description"`
	Variants    []Variant `json:"variants" gorm:"foreignKey:ProductID"`
}
