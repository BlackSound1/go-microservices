package data

import (
	"database/sql"
	"time"
)

type PostgresTestRepository struct {
	Conn *sql.DB
}

func NewPostgresTestRepository(db *sql.DB) *PostgresTestRepository {
	return &PostgresTestRepository{Conn: db}
}

// GetAll gets all users from the database.
func (u *PostgresTestRepository) GetAll() ([]*User, error) {
	users := []*User{}
	return users, nil
}

// GetByEmail gets a user by email.
func (u *PostgresTestRepository) GetByEmail(email string) (*User, error) {
	user := User{
		ID:        1,
		FirstName: "First",
		LastName:  "Last",
		Email:     "me@me.me",
		Password:  "",
		Active:    1,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	return &user, nil
}

// GetByID gets a user by ID.
func (u *PostgresTestRepository) GetByID(id int) (*User, error) {
	user := User{
		ID:        1,
		FirstName: "First",
		LastName:  "Last",
		Email:     "me@me.me",
		Password:  "",
		Active:    1,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	return &user, nil
}

// Update updates the user in the database.
func (u *PostgresTestRepository) Update(user User) error {
	return nil
}

// DeleteByID removes the user from the database by their ID.
func (u *PostgresTestRepository) DeleteByID(id int) error {
	return nil
}

// Insert creates a new user in the database, and returns the ID of the newly created user.
//
// It uses bcrypt to hash the password before storing it in the database.
func (u *PostgresTestRepository) Insert(user User) (int, error) {
	return 1, nil
}

// ResetPassword updates the user's password in the database.
//
// It hashes the new password using bcrypt before storing it.
func (u *PostgresTestRepository) ResetPassword(password string, user User) error {
	return nil
}

// PasswordMatches checks whether the given plaintext password matches the user's stored password.
func (u *PostgresTestRepository) PasswordMatches(plainText string, user User) (bool, error) {
	return true, nil
}
