package user

import (
	"links-service/pkg/db"
)

type UserRepository struct {
	Database *db.Db
}

func NewUserRepository(db *db.Db) *UserRepository {
	return &UserRepository{
		Database: db,
	}
}

func (repo *UserRepository) FindByEmail(email string) (*User, error) {
	var user User
	res := repo.Database.DB.First(&user, "email = ?", email)
	if res.Error != nil {
		return nil, res.Error
	}
	return &user, nil
}

func (repo *UserRepository) Create(user *User) (*User, error) {
	res := repo.Database.DB.Create(user)
	if res.Error != nil {
		return nil, res.Error
	}
	return user, nil
}
