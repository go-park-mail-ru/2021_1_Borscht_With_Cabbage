package http

import (
	"encoding/json"
	"github.com/borscht/backend/internal/models"
	"github.com/borscht/backend/internal/restaurantAdmin/mocks"
	"github.com/borscht/backend/utils/errors"
	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestNewSectionHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	sectionUsecaseMock := mocks.NewMockAdminSectionUsecase(ctrl)
	sectionHandler := NewSectionHandler(sectionUsecaseMock)
	if sectionHandler == nil {
		t.Errorf("incorrect result")
		return
	}
}

func TestSectionHandler_AddSection(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	sectionUsecaseMock := mocks.NewMockAdminSectionUsecase(ctrl)
	sectionHandler := NewSectionHandler(sectionUsecaseMock)

	newSection := models.Section{
		Name:       "yum",
		Restaurant: 1,
	}
	addedSection := newSection
	addedSection.ID = 1
	requestJSON := `{"name":"yum","restaurant":1}`

	e := echo.New()
	req := httptest.NewRequest(http.MethodPut, "/restaurant/section", strings.NewReader(requestJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	ctx := models.GetContext(c)

	sectionUsecaseMock.EXPECT().AddSection(ctx, newSection).Return(&addedSection, nil)

	err := sectionHandler.AddSection(c)
	if err != nil {
		t.Errorf("incorrect result")
		return
	}
}

func TestSectionHandler_AddSection_BindError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	sectionUsecaseMock := mocks.NewMockAdminSectionUsecase(ctrl)
	sectionHandler := NewSectionHandler(sectionUsecaseMock)

	requestJSON := `{"name""yum","restaurant":1}`

	e := echo.New()
	req := httptest.NewRequest(http.MethodPut, "/restaurant/section", strings.NewReader(requestJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := sectionHandler.AddSection(c)
	b := errors.SendError{}
	respCode := rec.Body.Bytes()
	err = json.Unmarshal(respCode, &b)
	if err != nil {
		t.Errorf("incorrect result")
		return
	}

	if b.Code == 200 {
		t.Errorf("incorrect result")
		return
	}
}

func TestSectionHandler_DeleteSection(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	sectionUsecaseMock := mocks.NewMockAdminSectionUsecase(ctrl)
	sectionHandler := NewSectionHandler(sectionUsecaseMock)

	newSection := models.Section{
		Name: "yum",
	}
	requestJSON := `{"name":"yum"}`
	deleted := models.DeleteSuccess{ID: 1}

	e := echo.New()
	req := httptest.NewRequest(http.MethodDelete, "/restaurant/section", strings.NewReader(requestJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	ctx := models.GetContext(c)

	sectionUsecaseMock.EXPECT().DeleteSection(ctx, newSection.ID).Return(&deleted, nil)

	err := sectionHandler.DeleteSection(c)
	if err != nil {
		t.Errorf("incorrect result")
		return
	}
}

func TestSectionHandler_DeleteSection_BindError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	sectionUsecaseMock := mocks.NewMockAdminSectionUsecase(ctrl)
	sectionHandler := NewSectionHandler(sectionUsecaseMock)

	requestJSON := `{"name""yum"}`

	e := echo.New()
	req := httptest.NewRequest(http.MethodDelete, "/restaurant/section", strings.NewReader(requestJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := sectionHandler.DeleteSection(c)
	b := errors.SendError{}
	respCode := rec.Body.Bytes()
	err = json.Unmarshal(respCode, &b)
	if err != nil {
		t.Errorf("incorrect result")
		return
	}

	if b.Code == 200 {
		t.Errorf("incorrect result")
		return
	}
}

func TestSectionHandler_UpdateSection(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	sectionUsecaseMock := mocks.NewMockAdminSectionUsecase(ctrl)
	sectionHandler := NewSectionHandler(sectionUsecaseMock)

	newSection := models.Section{
		Name: "yum",
	}
	requestJSON := `{"name":"yum"}`
	updatedSection := newSection
	updatedSection.ID = 1

	e := echo.New()
	req := httptest.NewRequest(http.MethodDelete, "/restaurant/section", strings.NewReader(requestJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	ctx := models.GetContext(c)

	sectionUsecaseMock.EXPECT().UpdateSection(ctx, newSection).Return(&updatedSection, nil)
	err := sectionHandler.UpdateSection(c)
	if err != nil {
		t.Errorf("incorrect result")
		return
	}
}

func TestSectionHandler_UpdateSection_BindError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	sectionUsecaseMock := mocks.NewMockAdminSectionUsecase(ctrl)
	sectionHandler := NewSectionHandler(sectionUsecaseMock)

	requestJSON := `{"name""yum"}`

	e := echo.New()
	req := httptest.NewRequest(http.MethodDelete, "/restaurant/section", strings.NewReader(requestJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := sectionHandler.UpdateSection(c)
	b := errors.SendError{}
	respCode := rec.Body.Bytes()
	err = json.Unmarshal(respCode, &b)
	if err != nil {
		t.Errorf("incorrect result")
		return
	}

	if b.Code == 200 {
		t.Errorf("incorrect result")
		return
	}
}
