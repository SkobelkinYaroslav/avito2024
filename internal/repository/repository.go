package repository

import "database/sql"

type Repository struct {
	NoAuthRepo
	AuthOnlyRepo
	ModerationsOnlyRepo
}

func NewRepository(db *sql.DB) Repository {
	return Repository{
		NoAuthRepo:          NewNoAuthRepo(db),
		AuthOnlyRepo:        NewAuthOnlyRepo(db),
		ModerationsOnlyRepo: NewModerationsOnlyRepo(db),
	}
}
