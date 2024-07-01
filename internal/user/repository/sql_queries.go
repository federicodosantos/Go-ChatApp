package repository

const (
	createQuery = `INSERT INTO users(id, name, email, photo_link, created_at, updated_at)
	VALUES($1, $2, $3, $4, $5, $6)
	RETURNING *`
)