package internal

import (
	"context"
	"database/sql"

	model "github.com/abhiii71/orderStream/account/models"
)

type AccountRepository interface {
	Close() error
	PutAccount(ctx context.Context, a model.Account) (*model.Account, error)
	GetAccountByEmail(ctx context.Context, email string) (*model.Account, error)
}

type repo struct {
	db *sql.DB
}

func NewAccountRepository(db *sql.DB) AccountRepository {
	return &repo{db: db}
}

func (r *repo) Close() error {
	return r.db.Close()
}

func (r *repo) PutAccount(ctx context.Context, a model.Account) (*model.Account, error) {
	query := `Insert into accounts (name, email, password) VALUES($1, $2, $3)`

	_, err := r.db.Exec(query, a.Name, a.Email, a.Password)
	if err != nil {
		return nil, err
	}
	return &a, nil
}

func (r *repo) GetAccountByEmail(ctx context.Context, email string) (*model.Account, error) {
	var account model.Account
	query := `Select id, name, email, password FROM accounts where email=$1`

	err := r.db.QueryRowContext(ctx, query, email).Scan(&account.ID, &account.Name, &account.Email, &account.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil

		} else {
			return nil, err
		}
	}

	return &account, nil
}
