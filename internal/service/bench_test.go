package service

import (
	models "L0/internal/model"
	"testing"
	"time"

	"L0/internal/cache"
)

type mockRepo struct{}

func (m *mockRepo) SaveOrder(*models.Order) error              { return nil }
func (m *mockRepo) GetOrderByID(string) (*models.Order, error) { return nil, nil }
func (m *mockRepo) LoadAllOrders() ([]*models.Order, error)    { return nil, nil }
func (m *mockRepo) OrderExists(string) (bool, error)           { return false, nil }
func (m *mockRepo) Ping() error                                { return nil }
func (m *mockRepo) Close() error                               { return nil }

func generateValidOrder() *models.Order {
	return &models.Order{
		OrderUID:        "bench-uid",
		TrackNumber:     "TRACK123",
		Entry:           "entry",
		Locale:          "en_US",
		CustomerID:      "cust-1",
		DeliveryService: "DHL",
		DateCreated:     time.Now(),
		Delivery: models.Delivery{
			Name: "John", Phone: "+1234567890", Zip: "123456",
			City: "NY", Address: "Street 1", Region: "NY",
			Email: "john@example.com",
		},
		Payment: models.Payment{
			Transaction: "txn-1", Currency: "USD", Provider: "payme",
			Amount: 100.0, PaymentDT: 1234567890, Bank: "bank",
			DeliveryCost: 10.0, GoodsTotal: 90.0, CustomFee: 0,
		},
		Items: []models.Item{
			{ChrtID: 1, TrackNumber: "T1", Price: 50.0, Rid: "r1",
				Name: "Item", Sale: 0, Size: "M", TotalPrice: 50.0,
				NmID: 1, Brand: "A", Status: 1},
		},
	}
}

func BenchmarkValidateOrder(b *testing.B) {
	svc := NewOrderService(&mockRepo{}, cache.New(1e6, 10*time.Minute))
	order := generateValidOrder()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = svc.validateOrder(order)
	}
}

func BenchmarkCreateOrder(b *testing.B) {
	svc := NewOrderService(&mockRepo{}, cache.New(1e6, 10*time.Minute))
	order := generateValidOrder()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = svc.CreateOrder(order)
	}
}
