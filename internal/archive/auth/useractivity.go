package auth

import (
	"fmt"
	"log"
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
