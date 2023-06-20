package auth

import (
	"errors"
	"log"
	"time"

	"github.com/ARTM2000/archive1/internal/archive/xerrors"
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

func NewUserRepository(db *gorm.DB) UserRepository {
	return UserRepository{
		db: db,
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
			return nil, xerrors.ErrRecordNotFound
		}
		log.Default().Println("[Unhandled] error in find admin user.", dbResult.Error.Error())
		return nil, xerrors.ErrUnhandled
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
			return nil, xerrors.ErrDuplicateViolation
		}
		log.Default().Println("[Unhandled] error in create admin user.", dbResult.Error.Error())
		return nil, xerrors.ErrUnhandled
	}

	return &newAdminUser, nil
}

func (repo *UserRepository) CreateNewNonAdminUser(username string, email string, hashedPassword string) (*UserSchema, error) {
	var newNonAdminUser = UserSchema{
		Username:       username,
		Email:          email,
		HashedPassword: hashedPassword,
		IsAdmin:        false,
	}
	dbResult := repo.db.Model(&UserSchema{}).Create(&newNonAdminUser)

	if dbResult.Error != nil {
		if errors.Is(dbResult.Error, gorm.ErrDuplicatedKey) {
			log.Default().Println("error in create admin user.", dbResult.Error.Error())
			return nil, xerrors.ErrDuplicateViolation
		}
		log.Default().Println("[Unhandled] error in create admin user.", dbResult.Error.Error())
		return nil, xerrors.ErrUnhandled
	}

	return &newNonAdminUser, nil
}

func (repo *UserRepository) FindUserWithEmail(email string) (*UserSchema, error) {
	var user UserSchema
	dbResult := repo.db.Model(&UserSchema{}).Where(UserSchema{Email: email}).First(&user)

	if dbResult.Error != nil {
		if errors.Is(dbResult.Error, gorm.ErrRecordNotFound) {
			log.Default().Println("error in find user with email.", dbResult.Error.Error())
			return nil, xerrors.ErrRecordNotFound
		}
		log.Default().Println("[Unhandled] error in find user with email.", dbResult.Error.Error())
		return nil, xerrors.ErrUnhandled
	}

	return &user, nil
}

func (repo *UserRepository) FindUserWithEmailOrUsername(email string, username string) (*UserSchema, error) {
	var user UserSchema
	dbResult := repo.db.Model(&UserSchema{}).Where(UserSchema{Email: email}).Or(UserSchema{Username: username}).First(&user)

	if dbResult.Error != nil {
		if errors.Is(dbResult.Error, gorm.ErrRecordNotFound) {
			log.Default().Println("error in find user with email.", dbResult.Error.Error())
			return nil, xerrors.ErrRecordNotFound
		}
		log.Default().Println("[Unhandled] error in find user with email.", dbResult.Error.Error())
		return nil, xerrors.ErrUnhandled
	}

	return &user, nil
}

func (repo *UserRepository) FindUserWithId(id uint) (*UserSchema, error) {
	var user UserSchema
	dbResult := repo.db.Model(&UserSchema{}).Where(UserSchema{ID: id}).First(&user)

	if dbResult.Error != nil {
		if errors.Is(dbResult.Error, gorm.ErrRecordNotFound) {
			log.Default().Println("error in find user with email.", dbResult.Error.Error())
			return nil, xerrors.ErrRecordNotFound
		}
		log.Default().Println("[Unhandled] error in find user with email.", dbResult.Error.Error())
		return nil, xerrors.ErrUnhandled
	}

	return &user, nil
}
