package models

import (
	"time"

	"github.com/uptrace/bun"
)

type Order struct {
	bun.BaseModel `bun:"table:orders,alias:o"`

	ID        string    `bun:"id,pk,type:uuid,default:gen_random_uuid()" json:"id"`
	UserID    string    `bun:"user_id,notnull" json:"user_id"`
	Status    string    `bun:"status,notnull,default:'pending'" json:"status"`
	Total     float64   `bun:"total,notnull" json:"total"`
	CreatedAt time.Time `bun:"created_at,notnull,default:now()" json:"created_at"`
	UpdatedAt time.Time `bun:"updated_at,notnull,default:now()" json:"updated_at"`
}

type OrderItem struct {
	bun.BaseModel `bun:"table:order_items,alias:oi"`

	ID        string  `bun:"id,pk,type:uuid,default:gen_random_uuid()" json:"id"`
	OrderID   string  `bun:"order_id,notnull" json:"order_id"`
	ProductID string  `bun:"product_id,notnull" json:"product_id"`
	Quantity  int     `bun:"quantity,notnull" json:"quantity"`
	Price     float64 `bun:"price,notnull" json:"price"`
}
