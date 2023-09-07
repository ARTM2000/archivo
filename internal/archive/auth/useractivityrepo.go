package auth

import (
	"log"
	"time"

	"github.com/ARTM2000/archive1/internal/archive/xerrors"
	"gorm.io/gorm"
)

type UserActivity struct {
	ID        uint `gorm:"primaryKey;unique" json:"id"`
	UserID    uint
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
