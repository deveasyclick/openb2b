package org

import (
	"github.com/deveasyclick/openb2b/internal/model"
	"gorm.io/gorm"
)

type WorkspaceRepository interface {
	Create(workspace *model.Org) error
	Update(workspace *model.Org) error
	Delete(ID uint) error
	FindByID(ID uint) (*model.Org, error)
}

type workspaceRepository struct {
	db *gorm.DB
}

func (r *workspaceRepository) Create(workspace *model.Org) error {
	return r.db.Create(workspace).Error
}

func (r *workspaceRepository) Update(workspace *model.Org) error {
	return r.db.Save(workspace).Error
}

func (r *workspaceRepository) Delete(ID uint) error {
	return r.db.Delete(&model.Org{}, ID).Error
}

func (r *workspaceRepository) FindByID(ID uint) (*model.Org, error) {
	var workspace model.Org
	err := r.db.First(&workspace, ID).Error
	if err != nil {
		return nil, err
	}
	return &workspace, nil
}

func NewWorkspaceRepository(db *gorm.DB) WorkspaceRepository {
	return &workspaceRepository{db: db}
}
