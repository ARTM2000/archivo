package database

import (
	"errors"
	"log"
	"time"

	"gorm.io/gorm"
)

type UserSchema struct {
	ID             uint           `gorm:"primaryKey;unique" json:"id"`
	Username       string         `gorm:"type:string;not null;unique" json:"username"`
	Email          string         `gorm:"type:string;not null;unique" json:"email"`
	HashedPassword string         `gorm:"type:string;not null" json:"-"`
	IsAdmin        bool           `gorm:"type:bool;not null" json:"is_admin"`
	CreatedAt      time.Time      `gorm:"autoUpdateTime:milli" json:"created_at"`
	UpdatedAt      time.Time      `gorm:"autoUpdateTime:milli" json:"updated_at"`
	DeletedAt      gorm.DeletedAt `json:"-"`
}

func (dbm *Manager) NewUserRepository() UserRepository {
	return UserRepository{
		db: dbm.db,
	}
}

type UserRepository struct {
	db *gorm.DB
}

func (repo *UserRepository) FindAdminUser() (*UserSchema, error) {
	var adminUser UserSchema
	dbResult := repo.db.Model(&UserSchema{}).Where(UserSchema{IsAdmin: true}).First(&adminUser)
	if dbResult.Error != nil {
		if errors.Is(dbResult.Error, gorm.ErrRecordNotFound) {
			log.Default().Println("error in find admin user.", dbResult.Error.Error())
			return nil, ErrRecordNotFound
		}
		log.Default().Println("[Unhandled] error in find admin user.", dbResult.Error.Error())
		return nil, ErrUnhandled
	}

	return &adminUser, nil
}

func (repo *UserRepository) CreateNewAdminUser(username string, email string, hashedPassword string) (*UserSchema, error) {
	var newAdminUser = UserSchema{
		Username:       username,
		Email:          email,
		HashedPassword: hashedPassword,
		IsAdmin:        true,
	}
	dbResult := repo.db.Model(&UserSchema{}).Create(&newAdminUser)

	if dbResult.Error != nil {
		if errors.Is(dbResult.Error, gorm.ErrDuplicatedKey) {
			log.Default().Println("error in create admin user.", dbResult.Error.Error())
			return nil, ErrDuplicateViolation
		}
		log.Default().Println("[Unhandled] error in create admin user.", dbResult.Error.Error())
		return nil, ErrUnhandled
	}

	return &newAdminUser, nil
}

func (repo *UserRepository) FindUserWithEmail(email string) (*UserSchema, error) {
	var user UserSchema
	dbResult := repo.db.Model(&UserSchema{}).Where(UserSchema{Email: email}).First(&user)

	if dbResult.Error != nil {
		if errors.Is(dbResult.Error, gorm.ErrRecordNotFound) {
			log.Default().Println("error in find user with email.", dbResult.Error.Error())
			return nil, ErrRecordNotFound
		}
		log.Default().Println("[Unhandled] error in find user with email.", dbResult.Error.Error())
		return nil, ErrUnhandled
	}

	return &user, nil
}
