package user

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID uuid.UUID `db:"id" json:"id"`
	Name string `db:"name" json:"name"`
	Email string `db:"email" json:"email"`
	Password string `db:"password" json:"password"`
	Photo_Link sql.NullString `db:"photo_link" json:"photo_link"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}