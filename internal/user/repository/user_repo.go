package repository

import (
	"fmt"

	"github.com/federicodosantos/Go-ChatApp/internal/user"
	"github.com/jmoiron/sqlx"
)

type UserRepoItf interface {
	CreateUser(user *user.User) error
}

type UserRepo struct {
	db *sqlx.DB
}


func NewUserRepo(db *sqlx.DB) UserRepoItf {
	return UserRepo{db: db}
}

// CreateUser implements UserRepoItf.
func (u UserRepo) CreateUser(user *user.User) error {
	result, err := u.db.Exec(createQuery, user.ID,
		 user.Name, user.Email, user.Photo_Link, user.CreatedAt, user.UpdatedAt)
	if err != nil {
		return err
	}
	
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows != 1 {
		return fmt.Errorf("expected single row affected, got %d rows affected", rows)
	}

	return nil
}