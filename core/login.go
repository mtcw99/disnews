package core

import (
	"golang.org/x/crypto/bcrypt"
)

type Login struct {
	Username string
	Hashpass []byte // Hashed Password
}

func (l *Login) Validate(plainpass string) bool {
	err := bcrypt.CompareHashAndPassword(l.Hashpass, []byte(plainpass))
	return err == nil
}

func LoginCreate(username string, plainpass string) (Login, error) {
	var login Login
	var err error
	login.Username = username
	login.Hashpass, err = bcrypt.GenerateFromPassword([]byte(plainpass), bcrypt.DefaultCost)
	if err != nil {
		return Login{}, err
	}
	return login, nil
}
