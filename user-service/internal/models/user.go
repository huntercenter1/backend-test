package models

import (
	"time"

	"github.com/uptrace/bun"
)

type User struct {
	bun.BaseModel `bun:"table:users,alias:u"`

	ID           string    `bun:"id,pk,type:uuid,default:gen_random_uuid()" json:"id"`
	Username     string    `bun:"username,unique,notnull" json:"username"`
	Email        string    `bun:"email,unique,notnull" json:"email"`
	PasswordHash string    `bun:"password_hash,notnull" json:"-"`
	CreatedAt    time.Time `bun:"created_at,notnull,default:now()" json:"created_at"`
	UpdatedAt    time.Time `bun:"updated_at,notnull,default:now()" json:"updated_at"`
}
