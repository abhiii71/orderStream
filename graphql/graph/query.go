package graph

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/abhiii71/orderStream/graphql/generated"
	"github.com/abhiii71/orderStream/graphql/models"
	"github.com/abhiii71/orderStream/graphql/utils"
	"github.com/abhiii71/orderStream/pkg/auth"
)

type queryResolver struct {
	server *Server
}

func (r *queryResolver) Accounts(ctx context.Context, pagination *generated.PaginationInput, id *int) ([]*models.Account, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	if id != nil {
		res, err := r.server.accountClient.GetAccount(ctx, uint64(*id))
		if err != nil {
			log.Println(err)
			return nil, err
		}
		return []*models.Account{{
			ID:    uint64(res.ID),
			Name:  res.Name,
			Email: res.Email,
		}}, nil
	}

	skip, take := uint64(0), uint64(0)
	if pagination != nil {
		skip, take = utils.Bounds(pagination)
	}
	accountList, err := r.server.accountClient.GetAccounts(ctx, skip, take)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	var accounts []*models.Account
	for _, account := range accountList {
		account := &models.Account{
			ID:    account.ID,
			Name:  account.Name,
			Email: account.Email,
		}

		accounts = append(accounts, account)
	}
	return accounts, nil
}

func (r *queryResolver) Product(ctx context.Context, pagination *generated.PaginationInput, query, id *string, viewedProductIds []*string, byAccountId *bool) ([]*generated.Product, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	// Get single

	if id != nil {
		res, err := r.server.productClient.GetProduct(ctx, *id)
		if err != nil {
			log.Println(err)
			return nil, err
		}

		return []*generated.Product{{
			ID:          res.Id,
			Name:        res.Name,
			Description: res.Description,
			Price:       res.Price,
		}}, nil
	}
	skip, take := uint64(0), uint64(0)
	if pagination != nil {
		skip, take = utils.Bounds(pagination)
	}

	// Get recommendations
	if viewedProductIds != nil {
		productIds := make([]string, len(viewedProductIds))
		for i, id := range viewedProductIds {
			productIds[i] = *id
		}
		res, err := r.server.recommenderClient.GetRecommendationBasedOnViewed(ctx, productIds, skip, take)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		productList := res.GetRecommendedProducts()
		var products []*generated.Product
		for _, product := range productList {
			products = append(products, &generated.Product{
				ID:          product.Id,
				Name:        product.Name,
				Description: product.Description,
				Price:       product.Price,
			})
		}
		return products, nil
	}
	if byAccountId != nil && *byAccountId {
		accountId := auth.GetUserId(ctx, true)
		if accountId == "" {
			return nil, errors.New("unauthorized")
		}
		skip = 0
		take = 100
		res, err := r.server.recommenderClient.GetRecommendationForUser(ctx, accountId, skip, take)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		productList := res.GetRecommendedProducts()
		var products []*generated.Product
		for _, product := range productList {

			products = append(products, &generated.Product{
				ID:          product.Id,
				Name:        product.Name,
				Description: product.Description,
				Price:       product.Price,
			})
		}
		return products, nil
	}
	q := ""
	if query != nil {
		q = *query
	}

	productList, err := r.server.productClient.GetProducts(ctx, skip, take, nil, q)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	var products []*generated.Product
	for _, product := range productList {
		products = append(products, &generated.Product{
			ID:          product.Id,
			Name:        product.Name,
			Description: product.Description,
			Price:       product.Price,
		})
	}
	return products, nil
}
