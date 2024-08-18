package service

import (
	"avito2024/internal/domain"
)

func (s Service) CreateHouseService(house domain.House) (domain.House, error) {
	if house.Address == "" || house.Year < 0 {
		return house, domain.NewCustomError(domain.InvalidInputError())
	}
	createdHouse, err := s.CreateHouseRepo(house)
	if err != nil {
		return house, err
	}
	return createdHouse, nil
}

func (s Service) UpdateFlatService(fl domain.Flat) (domain.Flat, error) {
	if fl.ID <= 0 {
		return fl, domain.NewCustomError(domain.InvalidInputError())
	}
	updatedFlat, err := s.UpdateFlatRepo(fl)
	if err != nil {
		return fl, err
	}
	return updatedFlat, nil
}
