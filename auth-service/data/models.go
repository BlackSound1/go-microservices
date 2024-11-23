package data

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"time"

	"golang.org/x/crypto/bcrypt"
)

const DB_TIMEOUT = time.Second * 3

var db *sql.DB

// New creates a Models object with the given database pool.
func New(dbPool *sql.DB) Models {
	db = dbPool

	return Models{User: User{}}
}

// Models holds all the models we will use
type Models struct {
	User User
}

// User holds a User from the database
type User struct {
	ID        int       `json:"id"`
	Email     string    `json:"email"`
	FirstName string    `json:"first_name,omitempty"`
	LastName  string    `json:"last_name,omitempty"`
	Password  string    `json:"-"`
	Active    int       `json:"active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// GetAll gets all users from the database.
func (u *User) GetAll() ([]*User, error) {

	// To avoid long queries
	ctx, cancel := context.WithTimeout(context.Background(), DB_TIMEOUT)
	defer cancel()

	query := `
		SELECT 
			id, email, first_name, last_name, password, user_active, created_at, updated_at
		FROM
			users
		ORDER BY
			last_name
	`

	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*User

	for rows.Next() {
		var user User

		err := rows.Scan(
			&user.ID,
			&user.Email,
			&user.FirstName,
			&user.LastName,
			&user.Password,
			&user.Active,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			log.Println("error scanning", err)
			return nil, err
		}
		users = append(users, &user)
	}
	return users, nil
}

// GetByEmail gets a user by email.
func (u *User) GetByEmail(email string) (*User, error) {

	// To avoid long queries
	ctx, cancel := context.WithTimeout(context.Background(), DB_TIMEOUT)
	defer cancel()

	query := `
		SELECT 
			id, email, first_name, last_name, password, user_active, created_at, updated_at
		FROM
			users
		WHERE
			email = $1
	`

	var user User
	row := db.QueryRowContext(ctx, query, email)

	err := row.Scan(
		&user.ID,
		&user.Email,
		&user.FirstName,
		&user.LastName,
		&user.Password,
		&user.Active,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// GetByID gets a user by ID.
func (u *User) GetByID(id int) (*User, error) {

	// To avoid long queries
	ctx, cancel := context.WithTimeout(context.Background(), DB_TIMEOUT)
	defer cancel()

	query := `
		SELECT
			id, email, first_name, last_name, password, user_active, created_at, updated_at
		FROM
			users
		WHERE
			id = $1
	`

	var user User
	row := db.QueryRowContext(ctx, query, id)

	err := row.Scan(
		&user.ID,
		&user.Email,
		&user.FirstName,
		&user.LastName,
		&user.Password,
		&user.Active,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// Update updates the user in the database.
func (u *User) Update() error {

	// To avoid long queries
	ctx, cancel := context.WithTimeout(context.Background(), DB_TIMEOUT)
	defer cancel()

	stmt := `
		UPDATE
			users
		SET
			email = $1,
			first_name = $2,
			last_name = $3,
			user_active = $4,
			updated_at = $5
		WHERE
			id = $6
	`

	_, err := db.ExecContext(
		ctx,
		stmt,
		u.Email,
		u.FirstName,
		u.LastName,
		u.Active,
		time.Now(),
		u.ID,
	)
	if err != nil {
		return err
	}

	return nil
}

// Delete removes the user from the database by their ID.
func (u *User) Delete() error {

	// To avoid long queries
	ctx, cancel := context.WithTimeout(context.Background(), DB_TIMEOUT)
	defer cancel()

	stmt := `
		DELETE FROM
			users
		WHERE
			id = $1
	`

	_, err := db.ExecContext(ctx, stmt, u.ID)
	if err != nil {
		return err
	}

	return nil
}

// DeleteByID removes the user from the database by their ID.
func (u *User) DeleteByID(id int) error {

	// To avoid long queries
	ctx, cancel := context.WithTimeout(context.Background(), DB_TIMEOUT)
	defer cancel()

	stmt := `
		DELETE FROM
			users
		WHERE
			id = $1
	`

	_, err := db.ExecContext(ctx, stmt, id)
	if err != nil {
		return err
	}

	return nil
}

// Insert creates a new user in the database, and returns the ID of the newly created user.
//
// It uses bcrypt to hash the password before storing it in the database.
func (u *User) Insert(user User) (int, error) {

	// To avoid long queries
	ctx, cancel := context.WithTimeout(context.Background(), DB_TIMEOUT)
	defer cancel()

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 12)
	if err != nil {
		return 0, err
	}

	var newID int

	stmt := `
		INSERT INTO
			users 
				(email, first_name, last_name, password, user_active, created_at, updated_at)
		VALUES
			($1, $2, $3, $4, $5, $6, $7)
		RETURNING id
	`

	err = db.QueryRowContext(
		ctx,
		stmt,
		user.Email,
		user.FirstName,
		user.LastName,
		hashedPassword,
		user.Active,
		time.Now(),
		time.Now(),
	).Scan(&newID)
	if err != nil {
		return 0, err
	}

	return newID, nil
}

// ResetPassword updates the user's password in the database.
//
// It hashes the new password using bcrypt before storing it.
func (u *User) ResetPassword(password string) error {

	// To avoid long queries
	ctx, cancel := context.WithTimeout(context.Background(), DB_TIMEOUT)
	defer cancel()

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}

	stmt := `
		UPDATE
			users
		SET
			password = $1
		WHERE
			id = $2
	`

	_, err = db.ExecContext(ctx, stmt, hashedPassword, u.ID)
	if err != nil {
		return err
	}

	return nil
}

// PasswordMatches checks whether the given plaintext password matches the user's stored password.
func (u *User) PasswordMatches(plainText string) (bool, error) {

	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(plainText))
	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			return false, nil
		default:
			return false, err
		}
	}
	return true, nil
}
