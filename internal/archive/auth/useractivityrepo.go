package auth

import (
	"log"
	"time"

	"github.com/ARTM2000/archivo/internal/archive/xerrors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type UserActivity struct {
	ID        uint      `gorm:"primaryKey;unique" json:"id"`
	UserID    uint      `json:"user_id"`
	Act       string    `gorm:"type:string;not null" json:"act"`
	CreatedAt time.Time `gorm:"autoUpdateTime:milli" json:"created_at"`
}

func NewUserActivityRepository(db *gorm.DB) UserActivityRepository {
	return UserActivityRepository{
		db,
	}
}

type UserActivityRepository struct {
	db *gorm.DB
}

func (uar *UserActivityRepository) SubmitNew(userId uint, act string) error {
	var newActivity = UserActivity{
		UserID: userId,
		Act:    act,
	}

	dbResult := uar.db.Model(&UserActivity{}).Create(&newActivity)
	if dbResult.Error != nil {
		log.Default().Printf("[Unhandled] error in creating new activity log, error: %+v", dbResult.Error)
		return xerrors.ErrUnhandled
	}

	return nil
}

func (uar *UserActivityRepository) SingleUserActivity(userId uint, option FindAllOption) (*[]UserActivity, int64, error) {
	var userActivities []UserActivity
	var DESC bool
	if option.SortOrder == "ASC" {
		DESC = false
	} else {
		DESC = true
	}

	dbResult := uar.db.Model(&UserActivity{}).Where(&UserActivity{UserID: userId}).Order(clause.OrderByColumn{Column: clause.Column{Name: option.SortBy}, Desc: DESC}).Offset(option.Start).Limit(option.End).Find(&userActivities)
	if dbResult.Error != nil {
		log.Default().Println("[Unhandled] error in finding user activity", dbResult.Error)
		return nil, 0, xerrors.ErrUnhandled
	}
	var total int64
	dbResult = uar.db.Model(&UserActivity{}).Where(&UserActivity{UserID: userId}).Order(clause.OrderByColumn{Column: clause.Column{Name: option.SortBy}, Desc: DESC}).Count(&total)
	if dbResult.Error != nil {
		log.Default().Println("[Unhandled] error in counting user activity", dbResult.Error)
		return nil, 0, xerrors.ErrUnhandled
	}
	return &userActivities, total, nil
}
