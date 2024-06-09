package serializer

import (
	"errors"
	"strings"
	"unicode"
)

type Token struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type RegisterRequest struct {
	Name     string `binding:"required,min=3,max=50" json:"name"`
	PhoneNo  string `binding:"required,min=11,max=11" json:"phone_no"`
	UserName string `binding:"max=50" json:"username"`
}

func (r RegisterRequest) PhoneNoValidate(phoneNo string) error {
	for _, item := range phoneNo {
		if !unicode.IsDigit(item) {
			return errors.New("bad_format")
		}
	}
	if !strings.HasPrefix(phoneNo, "09") {
		return errors.New("bad_format")
	}
	return nil
}
func (r RegisterRequest) UsernameValidate(username string) error {
	// handling ->  _ == 95 -> not first not last
	if username[0] == 95 || username[len(username)-1] == 95 {
		return errors.New("bad_format")
	}

	for index, item := range r.UserName {
		// handling -> _ == 95 -> not together
		if item == 95 && r.UserName[index+1] == 95 {
			return errors.New("bad_format")
		}
		// only digits and letter and _
		if !unicode.IsDigit(item) && !unicode.IsLetter(item) && item != 95 {
			return errors.New("bad_format")
		}

	}

	return nil
}

type LoginRequest struct {
	PhoneNo string `binding:"required,min=11,max=11" json:"phone_no"`
	Token   *int   `binding:"required,min=4,max=4" json:"token"`
}

func (l LoginRequest) PhoneNoValidate(phoneNo string) error {
	for _, item := range phoneNo {
		if !unicode.IsDigit(item) {
			return errors.New("bad_format")
		}
	}
	if !strings.HasPrefix(phoneNo, "09") {
		return errors.New("bad_format")
	}
	return nil
}
