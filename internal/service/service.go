package service

import "avito2024/internal/domain"

type NoAuthRepo interface {
	DummyLoginRepo(status string) (string, error)
	LoginRepo(user domain.User) (string, error)
	RegisterRepo(user domain.User) (string, error)
}
type AuthOnlyRepo interface {
	CheckUserRepo(tokenString string) (domain.User, error)
	GetHouseFlatsRepo(id int) ([]domain.Flat, error)
	CreateFlatRepo(flat domain.Flat) (domain.Flat, error)
}

type ModerationOnlyRepo interface {
	CreateHouseRepo(house domain.House) (domain.House, error)
	UpdateFlatRepo(fl domain.Flat) (domain.Flat, error)
}

type Repository interface {
	NoAuthRepo
	AuthOnlyRepo
	ModerationOnlyRepo
}

type Service struct {
	Repository
}

func (s Service) SubscribeHouseService(email string) error {
	//TODO implement me
	panic("implement me")
}

func NewService(repo Repository) Service {
	return Service{
		Repository: repo,
	}
}
