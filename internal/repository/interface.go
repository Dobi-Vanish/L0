package repository

import "L0/internal/model"

type Repository interface {
	SaveOrder(order *models.Order) error
	GetOrderByID(orderUID string) (*models.Order, error)
	LoadAllOrders() ([]*models.Order, error)
	OrderExists(orderUID string) (bool, error)
	Ping() error
	Close() error
}
