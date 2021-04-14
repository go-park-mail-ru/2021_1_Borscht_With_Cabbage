package repository

import (
	"context"
	"database/sql"
	"net/http"

	"github.com/borscht/backend/internal/models"
	"github.com/borscht/backend/internal/restaurantAdmin"
	"github.com/borscht/backend/utils/errors"
	"github.com/borscht/backend/utils/logger"
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
	dataToExistingCheck := models.CheckSectionExists{
		Name:         section.Name,
		RestaurantId: section.Restaurant,
	}
	err := s.checkExistingSection(ctx, dataToExistingCheck)
	if err != nil {
		return 0, err
	}

	var sid int

	err = s.DB.QueryRow("insert into sections (restaurant, name) values ($1, $2) returning sid",
		section.Restaurant, section.Name).Scan(&sid)
	if err != nil {
		failError := errors.FailServerError(err.Error())
		logger.RepoLevel().ErrorLog(ctx, failError)
		return 0, failError
	}

	return sid, nil
}

func (s sectionRepo) UpdateSection(ctx context.Context, section models.Section) error {
	dataToExistingCheck := models.CheckSectionExists{
		Name:         section.Name,
		RestaurantId: section.Restaurant,
		ID:           section.ID,
	}
	err := s.checkExistingSection(ctx, dataToExistingCheck)
	if err != nil {
		return err
	}

	_, err = s.DB.Exec(`update sections set name = $1
						where sid = $2`,
		section.Name, section.ID)
	if err != nil {
		failError := errors.FailServerError(err.Error())
		logger.RepoLevel().ErrorLog(ctx, failError)
		return failError
	}

	return nil
}

func (s sectionRepo) DeleteSection(ctx context.Context, sid int) error {
	_, err := s.DB.Exec("delete from sections where sid = $1", sid)
	if err != nil {
		failError := errors.FailServerError(err.Error())
		logger.RepoLevel().ErrorLog(ctx, failError)
		return failError
	}

	return nil
}

func (s sectionRepo) GetSection(ctx context.Context, sid int) (*models.Section, error) {
	DBsection, err := s.DB.Query("select sid, restaurant, name from sections where sid=$1", sid)
	if err != nil {
		return nil, errors.AuthorizationError("section not found")
	}

	section := new(models.Section)
	for DBsection.Next() {
		err = DBsection.Scan(
			&section.ID,
			&section.Restaurant,
			&section.Name,
		)
		if err != nil {
			failError := errors.FailServerError(err.Error())
			logger.RepoLevel().ErrorLog(ctx, failError)
			return nil, failError
		}
	}

	return section, nil
}

func (s sectionRepo) GetAllSections(ctx context.Context, idRestaurant int) ([]models.Section, error) {
	sectionsDB, err := s.DB.Query(`select sid, name from sections 
		where restaurant = $1`, idRestaurant)
	if err != nil {
		failError := errors.FailServerError(err.Error())
		logger.RepoLevel().ErrorLog(ctx, failError)
		return []models.Section{}, failError
	}

	var sections []models.Section
	for sectionsDB.Next() {
		section := new(models.Section)
		err = sectionsDB.Scan(
			&section.ID,
			&section.Name,
		)
		sections = append(sections, *section)
	}

	return sections, nil
}

func (s sectionRepo) checkExistingSection(ctx context.Context, sectionData models.CheckSectionExists) error {
	sections, err := s.DB.Query("select sid, name from sections where restaurant = $1", sectionData.RestaurantId)
	if err != nil {
		failError := errors.FailServerError(err.Error())
		logger.RepoLevel().ErrorLog(ctx, failError)
		return failError
	}
	for sections.Next() {
		nameSection := new(string)
		sidSection := new(int)
		sections.Scan(&sidSection, &nameSection)
		if *nameSection == sectionData.Name && *sidSection != sectionData.ID {
			return errors.NewCustomError(http.StatusBadRequest, "There is already such a section")
		}
	}

	return nil
}
