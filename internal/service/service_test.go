package service

import (
	"L0/internal/cache"
	"L0/internal/model"
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) SaveOrder(order *models.Order) error {
	args := m.Called(order)
	return args.Error(0)
}

func (m *MockRepository) GetOrderByID(orderUID string) (*models.Order, error) {
	args := m.Called(orderUID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Order), args.Error(1)
}

func (m *MockRepository) LoadAllOrders() ([]*models.Order, error) {
	args := m.Called()
	return args.Get(0).([]*models.Order), args.Error(1)
}

func (m *MockRepository) OrderExists(orderUID string) (bool, error) {
	args := m.Called(orderUID)
	return args.Bool(0), args.Error(1)
}

func (m *MockRepository) Ping() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockRepository) Close() error {
	args := m.Called()
	return args.Error(0)
}

func createValidOrder() models.Order {
	return models.Order{
		OrderUID:          "test123",
		TrackNumber:       "TRACK123",
		Entry:             "WBIL",
		Locale:            "en",
		InternalSignature: "",
		CustomerID:        "test_customer",
		DeliveryService:   "meest",
		Shardkey:          "9",
		SmID:              99,
		DateCreated:       time.Now(),
		OofShard:          "1",
		Delivery: models.Delivery{
			Name:    "Test Testov",
			Phone:   "+9720000000",
			Zip:     "2639809",
			City:    "Kiryat Mozkin",
			Address: "Ploshad Mira 15",
			Region:  "Kraiot",
			Email:   "test@gmail.com",
		},
		Payment: models.Payment{
			Transaction:  "test123",
			RequestID:    "",
			Currency:     "USD",
			Provider:     "wbpay",
			Amount:       1817,
			PaymentDT:    1637907727,
			Bank:         "alpha",
			DeliveryCost: 1500,
			GoodsTotal:   317,
			CustomFee:    0,
		},
		Items: []models.Item{
			{
				ChrtID:      9934930,
				TrackNumber: "TRACK123",
				Price:       453,
				Rid:         "ab4219087a764ae0btest",
				Name:        "Mascaras",
				Sale:        30,
				Size:        "0",
				TotalPrice:  317,
				NmID:        2389212,
				Brand:       "Vivienne Sabo",
				Status:      202,
			},
		},
	}
}

func TestOrderService_ProcessOrderMessage(t *testing.T) {
	mockRepo := new(MockRepository)
	cache := cache.New(1024*1024, time.Minute)

	service := NewOrderService(mockRepo, cache)

	order := createValidOrder()

	orderJSON, _ := json.Marshal(order)

	mockRepo.On("OrderExists", order.OrderUID).Return(false, nil)
	mockRepo.On("SaveOrder", mock.AnythingOfType("*models.Order")).Return(nil)

	err := service.ProcessOrderMessage(orderJSON)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestOrderService_ProcessOrderMessage_Invalid(t *testing.T) {
	mockRepo := new(MockRepository)
	cache := cache.New(1024*1024, time.Minute)

	service := NewOrderService(mockRepo, cache)

	invalidOrder := models.Order{
		TrackNumber: "TRACK123",
		CustomerID:  "customer123",
	}

	invalidOrderJSON, _ := json.Marshal(invalidOrder)

	err := service.ProcessOrderMessage(invalidOrderJSON)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "validation error")
	mockRepo.AssertNotCalled(t, "SaveOrder")
}

func TestOrderService_GetOrderByID(t *testing.T) {
	mockRepo := new(MockRepository)
	cache := cache.New(1024*1024, time.Minute)

	service := NewOrderService(mockRepo, cache)

	testOrder := createValidOrder()

	cache.Set(&testOrder)

	order, err := service.GetOrderByID("test123")
	assert.NoError(t, err)
	assert.Equal(t, testOrder.OrderUID, order.OrderUID)

	mockRepo.On("GetOrderByID", "test456").Return(&testOrder, nil)
	order, err = service.GetOrderByID("test456")
	assert.NoError(t, err)
	assert.Equal(t, testOrder.OrderUID, order.OrderUID)
	mockRepo.AssertExpectations(t)
}

func TestOrderService_CreateOrder(t *testing.T) {
	mockRepo := new(MockRepository)
	cache := cache.New(1024*1024, time.Minute)

	service := NewOrderService(mockRepo, cache)

	testOrder := createValidOrder()

	mockRepo.On("OrderExists", testOrder.OrderUID).Return(false, nil)
	mockRepo.On("SaveOrder", &testOrder).Return(nil)

	err := service.CreateOrder(&testOrder)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestOrderService_CreateOrder_AlreadyExists(t *testing.T) {
	mockRepo := new(MockRepository)
	cache := cache.New(1024*1024, time.Minute)

	service := NewOrderService(mockRepo, cache)

	testOrder := createValidOrder()

	mockRepo.On("OrderExists", testOrder.OrderUID).Return(true, nil)

	err := service.CreateOrder(&testOrder)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "already exists")
	mockRepo.AssertNotCalled(t, "SaveOrder")
}

func TestOrderService_RestoreCacheFromDB(t *testing.T) {
	mockRepo := new(MockRepository)
	cache := cache.New(1024*1024, time.Minute)

	service := NewOrderService(mockRepo, cache)

	testOrder := createValidOrder()
	orders := []*models.Order{&testOrder}

	mockRepo.On("LoadAllOrders").Return(orders, nil)

	err := service.RestoreCacheFromDB()
	assert.NoError(t, err)

	order, exists := cache.Get("test123")
	assert.True(t, exists)
	assert.Equal(t, testOrder.OrderUID, order.OrderUID)

	mockRepo.AssertExpectations(t)
}
