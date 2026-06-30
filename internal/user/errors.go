package user

import "errors"

var ErrUserAlreadyCreated = errors.New("error user already created")
var ErrInternalWhileCreateUser = errors.New("error internal while create user")
