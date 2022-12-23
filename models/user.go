package models

import (
	"errors"
	"regexp"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username string `json:"username"`
	Fullname string `json:"fullname"`
	Password string `json:"password"`
	ConfirmPassword string `json:"confirmpassword"`
}

func Register(db *gorm.DB, user *User) (err error) {
	err = db.Select("username", "fullname", "password").Create(user).Error
	if err != nil {
		return err
	}
	return nil
}

func GetUser(db *gorm.DB, user *User, username string) (err error) {
	err = db.Where("username = ?", username).First(user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
		return err
	}
	return nil
}

func Validate(user *User) (error) {
	invalidFormat:= ValidateEmaiil(user.Username)
	if invalidFormat != nil {
		return errors.New("Invalid Email Format")
	}

	if len(user.Fullname) < 3 {
		return errors.New("Name should be 2 characters or more")
	}

	invalidPassword:= ValidatePassword(user.Password, user.ConfirmPassword)
	if invalidPassword != nil {
		return invalidPassword
	}
	

	return nil
}

func ValidateEmaiil(email string) (error) {
	emailRegexp:= regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	if !emailRegexp.MatchString(email) {
		return errors.New("Invalid Email format")
	}
	return nil
}

func ValidatePassword(password string, confirmPassword string) (error) {
	if len(password) < 13 {
		return errors.New("Password should be at least 12 characters long")
	}

	if confirmPassword != "" {
		if password != confirmPassword {
			return errors.New("Confirmation Password doesn't match")
		}
	}
	return nil
}
