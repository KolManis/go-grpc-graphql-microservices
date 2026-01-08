package order

import "time"

type service interface {
	PutOrder
	GetordersForAccount
}

type Order struct {
	ID         string `json:""`
	CreatedAt  time.Time
	TotalPrice float64
	AccountID  string
	Products   []OrderedProduct
}
