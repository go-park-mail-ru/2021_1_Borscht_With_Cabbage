package auth

import (
	"context"
	"github.com/borscht/backend/internal/models"
)

type UserAuthRepo interface {
	Create(ctx context.Context, user models.User) (int, error)
	GetByUid(ctx context.Context, uid int) (models.User, error)
	GetByLogin(ctx context.Context, login string) (*models.User, error)
}

type RestaurantAuthRepo interface {
	CreateRestaurant(ctx context.Context, user models.RestaurantInfo) (int, error)
	GetByRid(ctx context.Context, rid int) (*models.RestaurantInfo, error)
	GetByLogin(ctx context.Context, login string) (*models.RestaurantInfo, error)
}

type SessionRepo interface {
	Create(ctx context.Context, sessionData models.SessionData) error
	Check(ctx context.Context, sessionToCheck string) (models.SessionInfo, bool, error)
	Delete(ctx context.Context, session string) error
}
