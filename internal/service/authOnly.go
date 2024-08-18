package service

import (
	"avito2024/internal/domain"
)

func (s Service) CheckUserService(tokenString string) (domain.User, error) {
	if tokenString == "" {
		return domain.User{}, domain.NewCustomError(domain.InvalidInputError())
	}
	var user domain.User
	user, err := s.CheckUserRepo(tokenString)
	if err != nil {
		return domain.User{}, err
	}
	return user, nil
}

func (s Service) GetHouseFlatsService(id int, user domain.User) ([]domain.Flat, error) {
	if id <= 0 {
		return nil, domain.NewCustomError(domain.InvalidInputError())
	}
	flats, err := s.GetHouseFlatsRepo(id)
	if err != nil {
		return nil, err
	}

	if user.UserType != "client" {
		return flats, nil
	}

	resultFlats := make([]domain.Flat, 0)
	for _, flat := range flats {
		if flat.Status == "approved" {
			resultFlats = append(resultFlats, flat)
		}
	}
	return resultFlats, nil

}

func (s Service) CreateFlatService(flat domain.Flat) (domain.Flat, error) {
	if flat.HouseID <= 0 || flat.Price <= 0 || flat.Rooms <= 0 {
		return domain.Flat{}, domain.NewCustomError(domain.InvalidInputError())
	}
	flat.Status = "created"
	createdFlat, err := s.CreateFlatRepo(flat)
	if err != nil {
		return domain.Flat{}, err
	}
	return createdFlat, nil
}
