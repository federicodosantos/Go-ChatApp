package user

import (
	"database/sql"
	"time"

)

type User struct {
	ID string `db:"id" json:"id"`
	Name string `db:"name" json:"name"`
	Email string `db:"email" json:"email"`
	Photo_Link sql.NullString `db:"photo_link" json:"photo_link"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}