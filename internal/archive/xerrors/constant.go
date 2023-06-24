package xerrors

import "errors"

var (
	ErrUnhandled                             = errors.New("unhandled error")
	ErrRecordNotFound                        = errors.New("record not found")
	ErrDuplicateViolation                    = errors.New("duplication violence")
	ErrAdminExist                            = errors.New("admin exists")
	ErrUserExist                             = errors.New("user exists")
	ErrEmailOrUsernameInUse                  = errors.New("email or username is used")
	ErrEmailOrPasswordIsIncorrect            = errors.New("email or password is incorrect")
	ErrUnauthorized                          = errors.New("unauthorized")
	ErrSourceServerWithThisNameExists        = errors.New("source server with this name exists")
	ErrUnableToCreateStoreDirectory          = errors.New("unable to create store directory")
	ErrStorePathExistButNotADirectory        = errors.New("store path exists but is not a directory")
	ErrFileRotateCountIsLowerThanPreviousOne = errors.New("file rotate count is lower than previous one")
)