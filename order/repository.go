package order

import (
	"context"
	"database/sql"

	"github.com/lib/pq"
)

type Repositoty interface {
	Close()
	PutOrder(ctx context.Context, o Order) error
	GetordersForAccount(ctx context.Context, accuntId string) ([]Order, error)
}

type postgresRepository struct {
	db *sql.DB
}

func (r *postgresRepository) Close() {
	r.db.Close()
}

func NewPostgresRepository(url string) (Repositoty, error) {
	db, err := sql.Open("postgres", url)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return &postgresRepository{db}, nil
}

func (r *postgresRepository) PutOrder(ctx context.Context, o Order) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}
		err = tx.Commit()
	}()
	_, err = tx.ExecContext(
		ctx,
		"INSERT INTO orders(id, created_at, account_id, total_price) VALUES ($1, $2, $3, $4)",
		o.ID,
		o.CreatedAt,
		o.AccountID,
		o.TotalPrice,
	)
	if err != nil {
		return err
	}

	stmt, err := tx.PrepareContext(ctx, pq.CopyIn("order_products", "order_id", "product_id", "quantity"))
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, p := range o.Products {
		_, err := stmt.ExecContext(ctx, o.ID, p.ID, p.quantity)
		if err != nil {
			return err
		}
	}

	_, err = stmt.ExecContext(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (r *postgresRepository) GetordersForAccount(ctx context.Context, accuntId string) ([]Order, error) {
	return
}
