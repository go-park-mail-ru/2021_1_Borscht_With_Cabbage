package repository

import (
	"context"
	"database/sql"

	"github.com/borscht/backend/internal/models"
	"github.com/borscht/backend/internal/restaurantAdmin"
)

type sectionRepo struct {
	DB *sql.DB
}

func NewSectionRepo(db *sql.DB) restaurantAdmin.AdminSectionRepo {
	return &sectionRepo{
		DB: db,
	}
}

func (s sectionRepo) AddSection(ctx context.Context, section models.Section) (int, error) {
	return 0, nil
}

func (s sectionRepo) UpdateSection(ctx context.Context, section models.Section) error {
	return nil
}

func (s sectionRepo) DeleteSection(ctx context.Context, sid int) error {
	return nil
}
