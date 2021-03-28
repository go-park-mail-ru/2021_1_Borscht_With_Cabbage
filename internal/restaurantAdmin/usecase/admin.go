package usecase

import (
	"github.com/borscht/backend/internal/models"
	"github.com/borscht/backend/internal/restaurantAdmin"
)

type adminUsecase struct {
	adminRepository restaurantAdmin.AdminRepo
}

func NewAdminUsecase(repo restaurantAdmin.AdminRepo) restaurantAdmin.AdminUsecase {
	return &adminUsecase{
		adminRepository: repo,
	}
}
func (a adminUsecase) Create(user models.Restaurant) (int32, error) {
	return a.adminRepository.Create(user)
}

func (a adminUsecase) CheckRestaurantExists(user models.RestaurantAuth) (models.Restaurant, error) {
	return a.adminRepository.CheckRestaurantExists(user)
}

func (a adminUsecase) GetByRid(rid int32) (models.Restaurant, error) {
	return a.adminRepository.GetByRid(rid)
}
