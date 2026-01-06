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
	// ExecContext выполняет запрос, не возвращая ни одной строки.
	// Аргументы предназначены для любых параметров-заполнителей в запросе.
	_, err := r.db.ExecContext(
		ctx,
		"INSERT INTO accounts(id,name) VALUES($1,$2)",
		a.ID, a.Name,
	)

	return err
}

func (r *postgrersRepository) GetAccountByID(ctx context.Context, id string) (*Account, error) {
	// QueryRowContext выполняет запрос, который, как ожидается, вернет не более одной строки.
	// QueryRowContext всегда возвращает ненулевое значение. Ошибки откладываются до
	// Вызывается метод сканирования [Row].
	// Если запрос не выбирает ни одной строки, [*Row.Scan] вернет [ErrNoRows].
	// В противном случае [*Row.Scan] сканирует первую выбранную строку и отбрасывает
	// остальное.
	a := &Account{}
	row := r.db.QueryRowContext(
		ctx,
		"SELECT id, name FROM accounts WHERE id = $1",
		id,
	)

	if err := row.Scan(&a.ID, &a.Name); err != nil {
		return nil, err
	}

	return a, nil
}

func (r *postgrersRepository) ListAccounts(ctx context.Context, skip uint64, take uint64) ([]Account, error) {
	rows, err := r.db.QueryContext(
		ctx,
		"SELECT id, name FROM accounts ORDER BY created_at DESC LIMIT $1 OFFSET $2",
		take,
		skip,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	accounts := make([]*Account, 0, take)
	for rows.Next() {
		a := &Account{}
		if err := rows.Scan(&a.ID, &a.Name); err != nil {
			return nil, err
		}
		accounts = append(accounts, *a)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return accounts, nil
}
