package logreg

import "errors"

var (
	ErrUsernameIsEmpty   error = errors.New("username is empty")
	ErrAlreadyRegistered error = errors.New("you already register -> try to login")
	ErrUserNotRegistered error = errors.New("user not registered -> try to register before login")
	ErrPasswordIncorrect error = errors.New("password incorrect")
)
