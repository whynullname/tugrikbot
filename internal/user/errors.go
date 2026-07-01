package user

import "errors"

var ErrUserAlreadyCreated = errors.New("error user already created")
var ErrInternalWhileCreateUser = errors.New("error internal while create user")
var ErrInternalWhileGetUserID = errors.New("internal error while get user id")
var ErrUserNotFound = errors.New("error user not found")
