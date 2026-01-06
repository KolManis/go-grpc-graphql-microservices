package account

import (
	"context"
	"database/sql"
)

type Repository interface {
	Close()
	PutAccount(ctx context.Context, a Account) error
	GetAccountByID(ctx context.Context, id string) (*Account, error)
	ListAccounts(ctx context.Context, skip uint64, take uint64) ([]Account, error)
}

type postgrersRepository struct {
	db *sql.DB
}

func NewPostgrersRepository(url string) (Repository, error) {
	db, err := sql.Open("postgres", url)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return &postgrersRepository{db}, nil
}

func (r *postgrersRepository) Close() {
	r.db.Close()
}

func (r *postgrersRepository) Ping() error {
	return r.db.Ping()
}

func (r *postgrersRepository) PutAccount(ctx context.Context, a Account) error {
	context

}

func (r *postgrersRepository) GetAccountByID(ctx context.Context, id string) (*Account, error) {

}

func (r *postgrersRepository) ListAccounts(ctx context.Context, skip uint64, take uint64) ([]Account, error) {
}
