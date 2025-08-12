package models

import (
	"time"

	"github.com/uptrace/bun"
)

type Product struct {
	bun.BaseModel `bun:"table:products,alias:p"`

	ID          string    `bun:"id,pk,type:uuid,default:gen_random_uuid()" json:"id"`
	Name        string    `bun:"name,notnull" json:"name"`
	Description string    `bun:"description" json:"description"`
	Price       float64   `bun:"price,notnull" json:"price"`
	Stock       int       `bun:"stock,notnull" json:"stock"`
	CreatedAt   time.Time `bun:"created_at,notnull,default:now()" json:"created_at"`
	UpdatedAt   time.Time `bun:"updated_at,notnull,default:now()" json:"updated_at"`
}
