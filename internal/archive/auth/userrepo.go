package auth

import (
	"errors"
	"log"
	"time"

	"github.com/ARTM2000/archive1/internal/archive/xerrors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type User struct {
	ID                    uint           `gorm:"primaryKey;unique" json:"id"`
	Username              string         `gorm:"type:string;not null;unique" json:"username"`
	Email                 string         `gorm:"type:string;not null;unique" json:"email"`
	HashedPassword        string         `gorm:"type:string;not null" json:"-"`
	IsAdmin               bool           `gorm:"type:bool;not null;default:false" json:"is_admin"`
	ChangeInitialPassword bool           `gorm:"type:bool;not null;default:true" json:"change_initial_password"`
	CreatedAt             time.Time      `gorm:"autoUpdateTime:milli" json:"created_at"`
	UpdatedAt             time.Time      `gorm:"autoUpdateTime:milli" json:"updated_at"`
	DeletedAt             gorm.DeletedAt `json:"-"`
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return UserRepository{
		db: db,
	}
}

type UserRepository struct {
	db *gorm.DB
}

type FindAllOption struct {
	SortBy    string
	SortOrder string
	Start     int
	End       int
}

func (repo *UserRepository) FindAdminUser() (*User, error) {
	var adminUser User
	dbResult := repo.db.Model(&User{}).Where(User{IsAdmin: true}).First(&adminUser)
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

func (repo *UserRepository) CreateNewAdminUser(username string, email string, hashedPassword string) (*User, error) {
	var newAdminUser = User{
		Username:              username,
		Email:                 email,
		HashedPassword:        hashedPassword,
		IsAdmin:               true,
		ChangeInitialPassword: false,
	}
	dbResult := repo.db.Model(&User{}).Create(&newAdminUser)

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

func (repo *UserRepository) CreateNewNonAdminUser(username string, email string, hashedPassword string) (*User, error) {
	var newNonAdminUser = User{
		Username:              username,
		Email:                 email,
		HashedPassword:        hashedPassword,
		IsAdmin:               false,
		ChangeInitialPassword: true,
	}
	dbResult := repo.db.Model(&User{}).Create(&newNonAdminUser)

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

func (repo *UserRepository) FindUserWithEmail(email string) (*User, error) {
	var user User
	dbResult := repo.db.Model(&User{}).Where(User{Email: email}).First(&user)

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

func (repo *UserRepository) FindUserWithEmailOrUsername(email string, username string) (*User, error) {
	var user User
	dbResult := repo.db.Model(&User{}).Where(User{Email: email}).Or(User{Username: username}).First(&user)

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

func (repo *UserRepository) FindUserWithId(id uint) (*User, error) {
	var user User
	dbResult := repo.db.Model(&User{}).Where(User{ID: id}).First(&user)

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

func (repo *UserRepository) FindAllUsers(option FindAllOption) (*[]User, int64, error) {
	var users []User
	var DESC bool
	if option.SortOrder == "ASC" {
		DESC = false
	} else {
		DESC = true
	}

	dbResult := repo.db.Model(&User{}).Order(clause.OrderByColumn{Column: clause.Column{Name: option.SortBy}, Desc: DESC}).Offset(option.Start).Limit(option.End).Find(&users)

	if dbResult.Error != nil {
		log.Default().Println("[Unhandled] error in find users.", dbResult.Error.Error())
		return nil, 0, xerrors.ErrUnhandled
	}

	return &users, dbResult.RowsAffected, nil
}

func (repo *UserRepository) ChangeUserPassword(id uint, newHashedPassword string) (*User, error) {
	user, err := repo.FindUserWithId(id)
	if err != nil {
		return nil, err
	}

	user.HashedPassword = newHashedPassword
	dbResult := repo.db.Save(user)
	if dbResult.Error != nil {
		log.Default().Printf("[Unhandled] error in changing user password, error: %+v", dbResult.Error)
		return nil, xerrors.ErrUnhandled
	}

	return user, nil
}
