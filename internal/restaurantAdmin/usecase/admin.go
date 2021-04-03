package usecase

import (
	"github.com/borscht/backend/config"
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
func (a adminUsecase) Create(restaurant models.Restaurant) (*models.Restaurant, error) {
	restaurant.Avatar = config.DefaultAvatar

	id, err := a.adminRepository.Create(restaurant)
	if err != nil {
		return nil, err
	}
	restaurant.ID = id
	return &restaurant, nil
}

func (a adminUsecase) CheckRestaurantExists(user models.RestaurantAuth) (*models.Restaurant, error) {
	return a.adminRepository.CheckRestaurantExists(user)
}

func (a adminUsecase) GetByRid(rid int) (models.Restaurant, error) {
	return a.adminRepository.GetByRid(rid)
}
