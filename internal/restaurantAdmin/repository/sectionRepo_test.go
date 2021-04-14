package repository

import (
	"context"
	"database/sql"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/borscht/backend/internal/models"
	"github.com/borscht/backend/utils/logger"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNewSectionRepo(t *testing.T) {

}

func TestSectionRepo_AddSection(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()
	sectionRepo := &sectionRepo{
		DB: db,
	}

	sectionInDB := sqlmock.NewRows([]string{"sid", "name"})
	sectionInDB.AddRow(2, "yum")

	sectionId := sqlmock.NewRows([]string{"sid"})
	sectionId.AddRow(1)

	section := models.Section{
		Name:       "yum1",
		Restaurant: 1,
	}

	mock.
		ExpectQuery("select sid, name from sections").
		WithArgs(1).
		WillReturnRows(sectionInDB)
	mock.
		ExpectQuery("insert into sections").
		WithArgs(section.Restaurant, section.Name).
		WillReturnRows(sectionId)

	ctx := new(context.Context)

	var id int
	id, err = sectionRepo.AddSection(*ctx, section)
	if err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}

	require.EqualValues(t, id, 1)
}

func TestSectionRepo_UpdateSection(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()
	sectionRepo := &sectionRepo{
		DB: db,
	}

	sectionInDB := sqlmock.NewRows([]string{"sid", "name"})
	sectionInDB.AddRow(2, "yum")

	section := models.Section{
		ID:         1,
		Name:       "yum1",
		Restaurant: 1,
	}

	mock.
		ExpectQuery("select sid, name from sections").
		WithArgs(1).
		WillReturnRows(sectionInDB)
	mock.
		ExpectExec("update sections set ").
		WithArgs(section.Name, section.ID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	ctx := new(context.Context)

	err = sectionRepo.UpdateSection(*ctx, section)
	if err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}
}

func TestSectionRepo_DeleteSection(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()
	sectionRepo := &sectionRepo{
		DB: db,
	}

	mock.
		ExpectExec("delete from sections where").
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	ctx := new(context.Context)

	err = sectionRepo.DeleteSection(*ctx, 1)
	if err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}
}

func TestSectionRepo_DeleteSectionError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()
	sectionRepo := &sectionRepo{
		DB: db,
	}

	mock.
		ExpectExec("delete from sections where").
		WithArgs(1).
		WillReturnError(sql.ErrNoRows)

	c := context.Background()
	ctx := context.WithValue(c, "request_id", 1)

	logger.InitLogger()
	err = sectionRepo.DeleteSection(ctx, 1)
	if err == nil {
		t.Errorf("unexpected err: %s", err)
		return
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}
}

func TestSectionRepo_GetSection(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()
	sectionRepo := &sectionRepo{
		DB: db,
	}

	sectionInDB := sqlmock.NewRows([]string{"sid", "restaurant", "name"})
	sectionInDB.AddRow(1, 1, "yum")

	mock.
		ExpectQuery("select sid, restaurant, name").
		WithArgs(1).
		WillReturnRows(sectionInDB)

	ctx := new(context.Context)

	sectionResult := new(models.Section)
	sectionResult, err = sectionRepo.GetSection(*ctx, 1)
	if err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}

	require.EqualValues(t, sectionResult.ID, 1)
}

func TestSectionRepo_GetAllSections(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()
	sectionRepo := &sectionRepo{
		DB: db,
	}

	sectionInDB := sqlmock.NewRows([]string{"sid", "name"})
	sectionInDB.AddRow(1, "section1")

	restaurantID := sqlmock.NewRows([]string{"restaurant"})
	restaurantID.AddRow(1)

	mock.
		ExpectQuery("select sid, name from").
		WithArgs(1).
		WillReturnRows(sectionInDB)

	ctx := new(context.Context)

	sectionResult := make([]models.Section, 0)
	sectionResult, err = sectionRepo.GetAllSections(*ctx, 1)
	if err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}

	require.EqualValues(t, sectionResult[0].ID, 1)
}
