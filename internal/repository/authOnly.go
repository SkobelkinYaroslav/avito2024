package repository

import (
	"avito2024/internal/domain"
	"database/sql"
)

type AuthOnlyRepo struct {
	db *sql.DB
}

func NewAuthOnlyRepo(db *sql.DB) AuthOnlyRepo {
	return AuthOnlyRepo{db: db}
}

func (a AuthOnlyRepo) CheckUserRepo(tokenString string) (domain.User, error) {
	var user domain.User
	query := "SELECT u.id, u.email, u.password, u.user_type FROM Users u JOIN Token t ON u.id = t.user_id WHERE t.token = $1"
	err := a.db.QueryRow(query, tokenString).Scan(&user.ID, &user.Email, &user.Password, &user.UserType)
	if err != nil {
		if err == sql.ErrNoRows {
			return domain.User{}, domain.NewCustomError(domain.NoFoundUserError())
		}
		return domain.User{}, domain.NewCustomError(domain.InternalError(err))
	}
	return user, nil
}

func (a AuthOnlyRepo) GetHouseFlatsRepo(id int) ([]domain.Flat, error) {
	var flats []domain.Flat
	query := "SELECT id, house_id, price, rooms, status FROM Flats WHERE house_id = $1"
	rows, err := a.db.Query(query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.NewCustomError(domain.InvalidInputError())
		}
		return nil, domain.NewCustomError(domain.InternalError(err))
	}
	defer rows.Close()

	for rows.Next() {
		var flat domain.Flat
		if err := rows.Scan(&flat.ID, &flat.HouseID, &flat.Price, &flat.Rooms, &flat.Status); err != nil {
			return nil, domain.NewCustomError(domain.InternalError(err))
		}
		flats = append(flats, flat)
	}
	if err := rows.Err(); err != nil {
		return nil, domain.NewCustomError(domain.InternalError(err))
	}
	return flats, nil
}

func (a AuthOnlyRepo) CreateFlatRepo(flat domain.Flat) (domain.Flat, error) {
	tx, err := a.db.Begin()
	if err != nil {
		return domain.Flat{}, domain.NewCustomError(domain.InternalError(err))
	}

	queryInsert := "INSERT INTO Flats (house_id, price, rooms, status) VALUES ($1, $2, $3, $4) RETURNING id"
	err = tx.QueryRow(queryInsert, flat.HouseID, flat.Price, flat.Rooms, flat.Status).Scan(&flat.ID)
	if err != nil {
		tx.Rollback()
		return domain.Flat{}, domain.NewCustomError(domain.InternalError(err))
	}

	queryUpdate := "UPDATE Houses SET updatedat = NOW() WHERE id = $1"
	_, err = tx.Exec(queryUpdate, flat.HouseID)
	if err != nil {
		tx.Rollback()
		return domain.Flat{}, domain.NewCustomError(domain.InternalError(err))
	}

	err = tx.Commit()
	if err != nil {
		return domain.Flat{}, domain.NewCustomError(domain.InternalError(err))
	}

	return flat, nil
}
