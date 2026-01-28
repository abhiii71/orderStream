package internal

import (
	"context"
	"errors"
	"log"

	"github.com/IBM/sarama"
	"github.com/abhiii71/orderStream/pkg/kafka"
	"github.com/abhiii71/orderStream/product/models"
)

type Service interface {
	GetProducer() sarama.AsyncProducer
	PostProduct(ctx context.Context, name, description string, price float64, accountId int) (*models.Product, error)
	GetProduct(ctx context.Context, id string) (*models.Product, error)
	GetProducts(ctx context.Context, skip, take uint64) ([]models.Product, error)
	GetProductsWithIds(ctx context.Context, ids []string) ([]models.Product, error)
	SearchProducts(ctx context.Context, query string, skip, take uint64) ([]models.Product, error)
	UpdateProduct(ctx context.Context, id, name, description string, price float64, accountId int) (*models.Product, error)
	DeleteProduct(ctx context.Context, productId string, accountId int) error
}

type productService struct {
	repo     Repository
	producer sarama.AsyncProducer
}

func NewProductService(repository Repository, producer sarama.AsyncProducer) Service {
	return &productService{repository, producer}
}

func (s *productService) GetProducer() sarama.AsyncProducer {
	return s.producer
}

func (s *productService) PostProduct(ctx context.Context, name, description string, price float64, accountId int) (*models.Product, error) {
	product := models.Product{
		Name:        name,
		Description: description,
		Price:       price,
		AccountId:   accountId,
	}

	err := s.repo.PutProduct(ctx, &product)
	if err != nil {
		return nil, err
	}

	go func() {
		err := kafka.SendMessageToRecommender(s, models.Event{
			Type: "product_created",
			Data: models.EventData{
				Id:          &product.Id,
				Name:        &product.Name,
				Description: &product.Description,
				Price:       &product.Price,
				AccountID:   &product.AccountId,
			},
		}, "product_events")
		if err != nil {
			log.Println("failed to send event to recommendation service: ", err)
		}
	}()
	return &product, nil
}

func (s *productService) GetProduct(ctx context.Context, id string) (*models.Product, error) {
	product, err := s.repo.GetProductsByID(ctx, id)
	if err != nil {
		return nil, err
	}

	go func() {
		err := kafka.SendMessageToRecommender(s, models.Event{
			Type: "product_retrieved",
			Data: models.EventData{
				Id:        &product.Id,
				AccountID: &product.AccountId,
			},
		}, "interaction_events")
		if err != nil {
			log.Println("failed to send event to recommendation service:", err)
		}
	}()

	return product, nil
}

func (s *productService) GetProducts(ctx context.Context, skip, take uint64) ([]models.Product, error) {
	products, err := s.repo.ListProducts(ctx, skip, take)
	if err != nil {
		return nil, err
	}

	// Send single products_listed event with all product IDs
	go func() {
		productIDs := make([]string, len(products))
		for i, product := range products {
			productIDs[i] = product.Id
		}

		err := kafka.SendMessageToRecommender(s, models.ProductsListedEvent{
			Type: "products_listed",
			Data: models.ProductsListedEventData{
				ProductIDs: productIDs,
				Count:      len(products),
			},
		}, "interaction_events")
		if err != nil {
			log.Println("failed to send event to recommendation service:", err)
		}
	}()

	return products, nil
}

func (s *productService) GetProductsWithIds(ctx context.Context, ids []string) ([]models.Product, error) {
	return s.repo.ListProductsWithIDs(ctx, ids)
}

func (s *productService) SearchProducts(ctx context.Context, query string, skip, take uint64) ([]models.Product, error) {
	return s.repo.SearchProducts(ctx, query, skip, take)
}

func (s *productService) UpdateProduct(ctx context.Context, id, name, description string, price float64, accountId int) (*models.Product, error) {
	product, err := s.repo.GetProductsByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if product.AccountId != accountId {
		return nil, errors.New("unauthorized")
	}

	updateProduct := &models.Product{
		Id:          id,
		Name:        name,
		Description: description,
		Price:       price,
		AccountId:   accountId,
	}

	err = s.repo.UpdateProduct(ctx, updateProduct)
	if err != nil {
		return nil, err
	}

	go func() {
		err := kafka.SendMessageToRecommender(s, models.Event{
			Type: "product_updated",
			Data: models.EventData{
				Id:          &updateProduct.Id,
				Name:        &updateProduct.Name,
				Description: &product.Description,
				Price:       &updateProduct.Price,
				AccountID:   &updateProduct.AccountId,
			},
		}, "product_events")
		if err != nil {
			log.Println("failed to send to recommendation service:", err)
		}
	}()

	return updateProduct, nil
}

func (s *productService) DeleteProduct(ctx context.Context, productId string, accountId int) error {
	product, err := s.repo.GetProductsByID(ctx, productId)
	if err != nil {
		return err
	}
	if product.AccountId != accountId {
		return errors.New("unauthorized")
	}

	go func() {
		err := kafka.SendMessageToRecommender(s, models.Event{
			Type: "product_deleted",
			Data: models.EventData{
				Id: &product.Id,
			},
		}, "product_events")
		if err != nil {
			log.Println("failed to send event to recommendation service:", err)
		}
	}()

	return s.repo.DeleteProduct(ctx, productId)
}
