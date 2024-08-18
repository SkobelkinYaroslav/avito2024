package repository

import (
	"avito2024/internal/domain"
	"avito2024/pkg"
	"database/sql"
	"fmt"
)

type NoAuthRepo struct {
	db *sql.DB
}

func NewNoAuthRepo(db *sql.DB) NoAuthRepo {
	return NoAuthRepo{db: db}
}

func (n NoAuthRepo) DummyLoginRepo(status string) (string, error) {
	tx, err := n.db.Begin()
	if err != nil {
		return "", domain.NewCustomError(domain.InternalError(err))
	}

	id := pkg.UUID()
	query := "INSERT INTO Users (id, email, password, user_type) VALUES ($1, $2, $3, $4)"
	token := pkg.UUID()
	dbData := fmt.Sprintf("dummy:%s", token)
	_, err = tx.Exec(query, id, dbData, dbData, status)
	if err != nil {
		tx.Rollback()
		return "", domain.NewCustomError(domain.InternalError(err))
	}

	query = "INSERT INTO Token (token, user_id) VALUES ($1, $2)"
	_, err = tx.Exec(query, token, id)
	if err != nil {
		tx.Rollback()
		return "", domain.NewCustomError(domain.InternalError(err))
	}

	err = tx.Commit()
	if err != nil {
		return "", domain.NewCustomError(domain.InternalError(err))
	}

	return token, nil
}

func (n NoAuthRepo) LoginRepo(user domain.User) (string, error) {
	var storedUser domain.User
	query := "SELECT t.token FROM Token t JOIN Users u ON u.id = t.user_id WHERE u.id = $1"
	err := n.db.QueryRow(query, user.ID).Scan(&storedUser.Token)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", domain.NewCustomError(domain.NoFoundUserError())
		}
		return "", domain.NewCustomError(domain.InternalError(err))
	}

	return storedUser.Token, nil
}

func (n NoAuthRepo) RegisterRepo(user domain.User) (string, error) {
	tx, err := n.db.Begin()
	if err != nil {
		return "", domain.NewCustomError(domain.InternalError(err))
	}

	id, token := pkg.UUID(), pkg.UUID()
	query := "INSERT INTO Users (id, email, password, user_type) VALUES ($1, $2, $3, $4)"
	_, err = tx.Exec(query, id, user.Email, user.Password, user.UserType)
	if err != nil {
		tx.Rollback()
		return "", domain.NewCustomError(domain.InternalError(err))
	}

	query = "INSERT INTO Token (token, user_id) VALUES ($1, $2)"
	_, err = tx.Exec(query, token, id)
	if err != nil {
		tx.Rollback()
		return "", domain.NewCustomError(domain.InternalError(err))
	}

	err = tx.Commit()
	if err != nil {
		return "", domain.NewCustomError(domain.InternalError(err))
	}

	return id, nil
}
