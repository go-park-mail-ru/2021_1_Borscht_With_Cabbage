package repository

import (
	"database/sql"
	"github.com/borscht/backend/internal/order"
)

type orderRepo struct {
	DB *sql.DB
}

func NewOrderRepo(db *sql.DB) order.OrderRepo {
	return &orderRepo{
		DB: db,
	}
}

func (o orderRepo) Create(session string, uid int) error {
	panic("implement me")
}

func (o orderRepo) GetUserOrder(uid int) {
	panic("implement me")
}
