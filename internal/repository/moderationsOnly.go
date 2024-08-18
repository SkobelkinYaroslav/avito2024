package repository

import (
	"avito2024/internal/domain"
	"database/sql"
)

type ModerationsOnlyRepo struct {
	db *sql.DB
}

func NewModerationsOnlyRepo(db *sql.DB) ModerationsOnlyRepo {
	return ModerationsOnlyRepo{db: db}
}

func (m ModerationsOnlyRepo) CreateHouseRepo(house domain.House) (domain.House, error) {
	query := "INSERT INTO Houses (address, year, developer) VALUES ($1, $2, $3) RETURNING id, createdat, updatedat"
	err := m.db.QueryRow(query, house.Address, house.Year, house.Developer).Scan(&house.ID, &house.CreatedAt, &house.UpdatedAt)
	if err != nil {
		return domain.House{}, domain.NewCustomError(domain.InternalError(err))
	}
	return house, nil
}

func (m ModerationsOnlyRepo) UpdateFlatRepo(fl domain.Flat) (domain.Flat, error) {
	query := "UPDATE flats SET status = 'on moderation' WHERE id = $1"
	_, err := m.db.Exec(query, fl.ID)
	if err != nil {
		return domain.Flat{}, domain.NewCustomError(domain.InternalError(err))
	}

	query = "UPDATE Flats SET status = $1 WHERE id = $2 RETURNING id, house_id, price, rooms, status"
	err = m.db.QueryRow(query, fl.Status, fl.ID).Scan(&fl.ID, &fl.HouseID, &fl.Price, &fl.Rooms, &fl.Status)
	if err != nil {
		if err == sql.ErrNoRows {
			return domain.Flat{}, domain.NewCustomError(domain.InvalidInputError())
		}
		return fl, domain.NewCustomError(domain.InternalError(err))
	}
	return fl, nil
}
