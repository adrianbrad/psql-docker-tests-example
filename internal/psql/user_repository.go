package psql

import (
	"context"
	"database/sql"
	"fmt"
)

// User represents a database user entity.
type User struct {
	ID    string
	Email string
}

// UserRepository stores the database connection on
// which CRUD actions are executed.
type UserRepository struct {
	db *sql.DB
}

// NewUserRepository creates an initialized UserRepository.
func NewUserRepository(db *sql.DB) *UserRepository {
	if db == nil {
		panic("nil db")
	}

	return &UserRepository{
		db: db,
	}
}

// AddUser inserts a User in the database.
func (u UserRepository) AddUser(
	ctx context.Context,
	user User,
) error {
	const insertStatement = `
	INSERT INTO users (user_id, email)
	VALUES ($1, $2)
	`

	_, err := u.db.ExecContext(
		ctx,
		insertStatement,
		user.ID,
		user.Email,
	)
	if err != nil {
		return fmt.Errorf("execute insert: %w", err)
	}

	return nil
}

// GetUser retrieves a user from the database.
func (u UserRepository) GetUser(
	ctx context.Context,
	id string,
) (User, error) {
	const queryStatement = `
	SELECT user_id, email 
		FROM users
		WHERE user_id=$1;
	`

	row := u.db.QueryRowContext(
		ctx,
		queryStatement,
		id,
	)

	var user User

	err := row.Scan(&user.ID, &user.Email)
	if err != nil {
		return User{}, fmt.Errorf("scan: %w", err)
	}

	return user, nil
}
