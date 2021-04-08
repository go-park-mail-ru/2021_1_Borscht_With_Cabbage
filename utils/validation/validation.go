package validation

import (
	"github.com/borscht/backend/utils/errors"
	"regexp"
)

var emailRegex = regexp.MustCompile("^([A-Za-z0-9_\\-.])+@([A-Za-z0-9_\\-.])+\\.([A-Za-z]{2,4})$")
var phoneNumberRegex = regexp.MustCompile("^[0-9]{11}$")

func ValidateEmail(email string) error {
	if len(email) == 0 {
		return errors.BadRequestError("email can't be empty")
	}

	if !emailRegex.MatchString(email) {
		return errors.BadRequestError("email is not valid")
	}

	return nil
}

func ValidatePhoneNumber(number string) error {
	if len(number) == 0 {
		return errors.BadRequestError("phone number can't be empty")
	}

	if !phoneNumberRegex.MatchString(number) {
		return errors.BadRequestError("phone number is not valid")
	}

	return nil
}

func ValidatePassword(password string) error {
	passwordLength := len(password)

	if passwordLength == 0 {
		return errors.BadRequestError("password can't be empty")
	}

	if passwordLength > 30 || passwordLength < 6 {
		return errors.BadRequestError("password must be 6-30 symbols")
	}

	return nil
}

func ValidateName(name string) error {
	if len(name) == 0 {
		return errors.BadRequestError("name can't be empty")
	}

	return nil
}

func ValidateLogin(login string) error {
	if !emailRegex.MatchString(login) && !phoneNumberRegex.MatchString(login) {
		return errors.BadRequestError("login is not valid")
	}

	return nil
}
