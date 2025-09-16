package service

import (
	"L0/internal/cache"
	"L0/internal/metrics"
	"L0/internal/model"
	"L0/internal/repository"
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"
)

type OrderService struct {
	repo  repository.Repository
	cache *cache.Cache
}

func NewOrderService(repo repository.Repository, cache *cache.Cache) *OrderService {
	return &OrderService{
		repo:  repo,
		cache: cache,
	}
}

func (s *OrderService) CreateOrder(order *models.Order) error {
	if err := s.validateOrder(order); err != nil {
		return fmt.Errorf("validation error: %v", err)
	}

	exists, err := s.repo.OrderExists(order.OrderUID)
	if err != nil {
		return fmt.Errorf("error checking order existence: %v", err)
	}

	if exists {
		return fmt.Errorf("order %s already exists", order.OrderUID)
	}

	if err := s.repo.SaveOrder(order); err != nil {
		metrics.DBErrors.Inc()
		return fmt.Errorf("error saving order: %v", err)
	}

	s.cache.Set(order)
	metrics.OrdersProcessed.Inc()

	return nil
}

func (s *OrderService) ProcessOrderMessage(message []byte) error {
	var order models.Order
	if err := json.Unmarshal(message, &order); err != nil {
		return fmt.Errorf("error parsing message: %v", err)
	}

	if err := s.validateOrder(&order); err != nil {
		return fmt.Errorf("validation error: %v", err)
	}

	exists, err := s.repo.OrderExists(order.OrderUID)
	if err != nil {
		return fmt.Errorf("error checking order existence: %v", err)
	}

	if exists {
		return fmt.Errorf("order %s already exists", order.OrderUID)
	}

	if err := s.repo.SaveOrder(&order); err != nil {
		metrics.DBErrors.Inc()
		return fmt.Errorf("error saving order: %v", err)
	}

	s.cache.Set(&order)
	metrics.OrdersProcessed.Inc()
	metrics.KafkaMessagesReceived.Inc()

	return nil
}

func (s *OrderService) GetOrderByID(orderUID string) (*models.Order, error) {
	if order, exists := s.cache.Get(orderUID); exists {
		metrics.OrdersFromCache.Inc()
		return order, nil
	}

	order, err := s.repo.GetOrderByID(orderUID)
	if err != nil {
		metrics.DBErrors.Inc()
		return nil, fmt.Errorf("error getting order from DB: %v", err)
	}

	if order != nil {
		s.cache.Set(order)
		metrics.OrdersFromDB.Inc()
	}

	return order, nil
}

func (s *OrderService) validateOrder(order *models.Order) error {
	if order.OrderUID == "" {
		return errors.New("order UID is required")
	}
	if len(order.OrderUID) > 255 {
		return errors.New("order UID must be less than 255 characters")
	}

	if order.TrackNumber == "" {
		return errors.New("track number is required")
	}
	if len(order.TrackNumber) > 255 {
		return errors.New("track number must be less than 255 characters")
	}

	if order.Entry == "" {
		return errors.New("entry is required")
	}
	if len(order.Entry) > 50 {
		return errors.New("entry must be less than 50 characters")
	}

	if order.Locale == "" {
		return errors.New("locale is required")
	}
	if len(order.Locale) > 10 {
		return errors.New("locale must be less than 10 characters")
	}

	if order.CustomerID == "" {
		return errors.New("customer ID is required")
	}
	if len(order.CustomerID) > 255 {
		return errors.New("customer ID must be less than 255 characters")
	}

	if order.DeliveryService == "" {
		return errors.New("delivery service is required")
	}
	if len(order.DeliveryService) > 100 {
		return errors.New("delivery service must be less than 100 characters")
	}

	if order.DateCreated.IsZero() {
		return errors.New("date created is required")
	}
	if order.DateCreated.After(time.Now().Add(24 * time.Hour)) {
		return errors.New("date created cannot be in the future")
	}

	if order.Delivery.Name == "" {
		return errors.New("delivery name is required")
	}
	if len(order.Delivery.Name) > 255 {
		return errors.New("delivery name must be less than 255 characters")
	}

	if order.Delivery.Phone == "" {
		return errors.New("delivery phone is required")
	}
	if !isValidPhone(order.Delivery.Phone) {
		return errors.New("delivery phone format is invalid")
	}

	if order.Delivery.Zip == "" {
		return errors.New("delivery zip is required")
	}
	if len(order.Delivery.Zip) > 50 {
		return errors.New("delivery zip must be less than 50 characters")
	}

	if order.Delivery.City == "" {
		return errors.New("delivery city is required")
	}
	if len(order.Delivery.City) > 255 {
		return errors.New("delivery city must be less than 255 characters")
	}

	if order.Delivery.Address == "" {
		return errors.New("delivery address is required")
	}
	if len(order.Delivery.Address) > 500 {
		return errors.New("delivery address must be less than 500 characters")
	}

	if order.Delivery.Region == "" {
		return errors.New("delivery region is required")
	}
	if len(order.Delivery.Region) > 255 {
		return errors.New("delivery region must be less than 255 characters")
	}

	if order.Delivery.Email == "" {
		return errors.New("delivery email is required")
	}
	if !isValidEmail(order.Delivery.Email) {
		return errors.New("delivery email format is invalid")
	}

	if order.Payment.Transaction == "" {
		return errors.New("payment transaction is required")
	}
	if len(order.Payment.Transaction) > 255 {
		return errors.New("payment transaction must be less than 255 characters")
	}

	if order.Payment.Currency == "" {
		return errors.New("payment currency is required")
	}
	if len(order.Payment.Currency) > 10 {
		return errors.New("payment currency must be less than 10 characters")
	}

	if order.Payment.Provider == "" {
		return errors.New("payment provider is required")
	}
	if len(order.Payment.Provider) > 100 {
		return errors.New("payment provider must be less than 100 characters")
	}

	if order.Payment.Amount <= 0 {
		return errors.New("payment amount must be positive")
	}

	if order.Payment.PaymentDT <= 0 {
		return errors.New("payment date is required")
	}

	if order.Payment.Bank == "" {
		return errors.New("payment bank is required")
	}
	if len(order.Payment.Bank) > 100 {
		return errors.New("payment bank must be less than 100 characters")
	}

	if order.Payment.DeliveryCost < 0 {
		return errors.New("payment delivery cost cannot be negative")
	}

	if order.Payment.GoodsTotal <= 0 {
		return errors.New("payment goods total must be positive")
	}

	if order.Payment.CustomFee < 0 {
		return errors.New("payment custom fee cannot be negative")
	}

	if len(order.Items) == 0 {
		return errors.New("order must have at least one item")
	}

	for i, item := range order.Items {
		if item.ChrtID <= 0 {
			return fmt.Errorf("item %d chrt_id must be positive", i)
		}

		if item.TrackNumber == "" {
			return fmt.Errorf("item %d track number is required", i)
		}
		if len(item.TrackNumber) > 255 {
			return fmt.Errorf("item %d track number must be less than 255 characters", i)
		}

		if item.Price <= 0 {
			return fmt.Errorf("item %d price must be positive", i)
		}

		if item.Rid == "" {
			return fmt.Errorf("item %d rid is required", i)
		}
		if len(item.Rid) > 255 {
			return fmt.Errorf("item %d rid must be less than 255 characters", i)
		}

		if item.Name == "" {
			return fmt.Errorf("item %d name is required", i)
		}
		if len(item.Name) > 255 {
			return fmt.Errorf("item %d name must be less than 255 characters", i)
		}

		if item.Sale < 0 {
			return fmt.Errorf("item %d sale cannot be negative", i)
		}

		if item.Size == "" {
			return fmt.Errorf("item %d size is required", i)
		}
		if len(item.Size) > 50 {
			return fmt.Errorf("item %d size must be less than 50 characters", i)
		}

		if item.TotalPrice <= 0 {
			return fmt.Errorf("item %d total price must be positive", i)
		}

		if item.NmID <= 0 {
			return fmt.Errorf("item %d nm_id must be positive", i)
		}

		if item.Brand == "" {
			return fmt.Errorf("item %d brand is required", i)
		}
		if len(item.Brand) > 255 {
			return fmt.Errorf("item %d brand must be less than 255 characters", i)
		}

		if item.Status < 0 {
			return fmt.Errorf("item %d status cannot be negative", i)
		}
	}

	return nil
}

func isValidEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,}$`)
	return emailRegex.MatchString(strings.ToLower(email))
}

func isValidPhone(phone string) bool {
	phoneRegex := regexp.MustCompile(`^\+?[0-9]{10,15}$`)
	return phoneRegex.MatchString(phone)
}

func (s *OrderService) RestoreCacheFromDB() error {
	orders, err := s.repo.LoadAllOrders()
	if err != nil {
		return fmt.Errorf("error loading orders from DB: %v", err)
	}

	s.cache.Restore(orders)
	return nil
}
