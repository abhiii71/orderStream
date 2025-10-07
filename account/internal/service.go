package internal

import (
	"context"
	"errors"

	model "github.com/abhiii71/orderStream/account/models"
	"github.com/abhiii71/orderStream/pkg/auth"
	"github.com/abhiii71/orderStream/pkg/crypt"
)

type AccountService interface {
	Register(ctx context.Context, name, email, password string) (string, error)
}

type service struct {
	repo AccountRepository
}

func Newservice(r AccountRepository) AccountService {
	return &service{r}
}

// func (s *service) Register(ctx context.Context, name, email, password string) (string, error) {
// 	_, err := s.repo.GetAccountByEmail(ctx, email)
// 	if err == nil {
// 		return "", errors.New("account already exists!")
// 	}

// 	// hased the password to save into db
// 	hashedPassword, err := crypt.HashPassword(password)
// 	if err != nil {
// 		return "", err
// 	}

// 	acc := model.Account{
// 		Name:     name,
// 		Email:    email,
// 		Password: hashedPassword,
// 	}

// 	// if email not exists already in db then create account
// 	account, err := s.repo.PutAccount(ctx, acc)
// 	if err != nil {
// 		return "", err
// 	}

// 	// generate token
// 	token, err := auth.GenerateToken(account.ID)
// 	if err != nil {
// 		return "", err
// 	}

// 	return token, nil
// }

func (s *service) Register(ctx context.Context, name, email, password string) (string, error) {
	account, err := s.repo.GetAccountByEmail(ctx, email)
	if err != nil {
		return "", err
	}
	if account != nil {
		return "", errors.New("account already exists!")
	}

	// hash password
	hashedPassword, err := crypt.HashPassword(password)
	if err != nil {
		return "", err
	}

	acc := model.Account{
		Name:     name,
		Email:    email,
		Password: hashedPassword,
	}

	account, err = s.repo.PutAccount(ctx, acc)
	if err != nil {
		return "", err
	}

	// generate token
	token, err := auth.GenerateToken(account.ID)
	if err != nil {
		return "", err
	}

	return token, nil
}
