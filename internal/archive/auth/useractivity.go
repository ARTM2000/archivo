package auth

import (
	"fmt"
	"log"

	"github.com/ARTM2000/archivo/internal/archive/xerrors"
)

func NewUserActivityManager(userActivityRepo UserActivityRepository) userActivityManager {
	return userActivityManager{
		userActivityRepo,
	}
}

type userActivityManager struct {
	userActivityRepo UserActivityRepository
}

func (uam *userActivityManager) SaveNewActivity(userId uint, method, route string) error {
	log.Default().Printf(
		"check user activity log for '%s:%s' for user %d",
		method,
		route,
		userId,
	)
	err := uam.userActivityRepo.SubmitNew(userId, fmt.Sprintf("%s:%s", method, route))
	return err
}

func (uam *userActivityManager) GetListForSingleUser(userId uint, option FindAllOption) (*[]UserActivity, int64, error) {
	activities, total, err := uam.userActivityRepo.SingleUserActivity(userId, option)

	if err != nil {
		log.Default().Printf("[Unhandled] error in getting list of user '%d' activity, error: %+v", userId, err)
		return nil, 0, xerrors.ErrUnhandled
	}

	return activities, total, nil
}
