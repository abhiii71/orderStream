package internal

import (
	"context"
	"log"
	"time"

	"github.com/IBM/sarama"
	"github.com/abhiii71/orderStream/order/models"
	"github.com/abhiii71/orderStream/pkg/kafka"
)

type Service interface {
	PostOrder(ctx context.Context, accountId uint64, totalPrice float64, products []*models.OrderedProduct) (*models.Order, error)
	GetOrdersForAccount(ctx context.Context, accountId uint64) ([]*models.Order, error)
	UpdateOrderPaymentStatus(ctx context.Context, orderId uint64, status string) error
	GetProducer() sarama.AsyncProducer
}

type orderService struct {
	repo     Repository
	producer sarama.AsyncProducer
}

func NewOrderService(repository Repository, producer sarama.AsyncProducer) Service {
	return &orderService{repository, producer}
}

func (s *orderService) PostOrder(ctx context.Context, accountId uint64, totalPrice float64, products []*models.OrderedProduct) (*models.Order, error) {
	order := models.Order{
		AccountID:  accountId,
		TotalPrice: totalPrice,
		Products:   products,
		CreatedAt:  time.Now().UTC(),
	}

	err := s.repo.PutOrder(ctx, &order)
	if err != nil {
		return nil, err
	}

	// send to recommendation service
	go func() {

		if err != nil {
			log.Println("failed to convert account ID to int:", err)
			return
		}

		for _, product := range products {
			err := kafka.SendMessageToRecommender(s, models.Event{
				Type: "purchase",
				Data: models.EventData{
					AccountId: int(accountId),
					ProductId: string(product.ID),
					// ProductId: strconv.FormatUint(uint64(product.ID), 10),
				},
			}, "interaction_events")
			if err != nil {
				log.Println("failed to send event to recommendation service: ", err)
			}
		}
	}()

	return &order, nil
}

func (s *orderService) GetOrdersForAccount(ctx context.Context, accountId uint64) ([]*models.Order, error) {
	return s.repo.GetOrdersForAccount(ctx, accountId)
}

func (s *orderService) UpdateOrderPaymentStatus(ctx context.Context, orderId uint64, status string) error {
	return s.repo.UpdateOrderPaymentStatus(ctx, orderId, status)
}

func (s *orderService) GetProducer() sarama.AsyncProducer {
	return s.producer
}
