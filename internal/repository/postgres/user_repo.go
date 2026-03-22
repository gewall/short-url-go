package repository

import (
	"database/sql"
	"errors"

	"github.com/gewall/short-url/internal/domain"
	"github.com/google/uuid"
)

type userRepository struct {
	db *sql.DB
}

func NewUserRepo(db *sql.DB) *userRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(user domain.User) (*domain.User, error) {
	var _user domain.User
	query := "INSERT INTO users (username, password) VALUES ($1, $2) RETURNING id,username,created_at"
	err := r.db.QueryRow(query, user.Username, user.Password).Scan(&_user.ID, &_user.Username, &_user.CreatedAt)
	if err != nil {

		return nil, err
	}

	return &_user, nil
}

func (r *userRepository) FindByUsername(username string) (*domain.User, error) {
	var _user domain.User
	query := "SELECT * FROM users WHERE username=$1"
	err := r.db.QueryRow(query, username).Scan(&_user.ID, &_user.Username, &_user.Password, &_user.CreatedAt)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &_user, nil
}

func (r *userRepository) FindByID(id uuid.UUID) (*domain.User, error) {
	var _user domain.User
	query := "SELECT id,username,created_at FROM users WHERE id=$1"
	err := r.db.QueryRow(query, id).Scan(&_user.ID, &_user.Username, &_user.CreatedAt)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &_user, nil
}

func (r *userRepository) FindAll() ([]*domain.User, error) {
	var users []*domain.User
	query := "SELECT id,username,created_at FROM users"
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		user := new(domain.User)

		err := rows.Scan(&user.ID, &user.Username, &user.CreatedAt)
		if err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func (r *userRepository) Delete(id uuid.UUID) error {
	query := "DELETE FROM users WHERE id=$1"
	_, err := r.db.Exec(query, id)

	if err != nil {
		return err
	}
	return nil
}
