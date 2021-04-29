package validation

import (
	"regexp"

	"github.com/borscht/backend/utils/errors"
)

var emailRegex = regexp.MustCompile("^([A-Za-z0-9_\\-.])+@([A-Za-z0-9_\\-.])+\\.([A-Za-z]{2,4})$")
var phoneNumberRegex = regexp.MustCompile("^[0-9]{11}$")

func ValidateEmail(email string) error {
	if len(email) == 0 {
		return errors.NewErrorWithMessage("email can't be empty")
	}

	if !emailRegex.MatchString(email) {
		return errors.NewErrorWithMessage("email is not valid")
	}

	return nil
}

func ValidatePhoneNumber(number string) error {
	if len(number) == 0 {
		return errors.NewErrorWithMessage("phone number can't be empty")
	}

	if !phoneNumberRegex.MatchString(number) {
		return errors.NewErrorWithMessage("phone number is not valid")
	}

	return nil
}

func ValidatePassword(password string) error {
	passwordLength := len(password)

	if passwordLength == 0 {
		return errors.NewErrorWithMessage("password can't be empty")
	}

	if passwordLength > 30 || passwordLength < 6 {
		return errors.NewErrorWithMessage("password must be 6-30 symbols")
	}

	return nil
}

func ValidateName(name string) error {
	if len(name) == 0 {
		return errors.NewErrorWithMessage("name can't be empty")
	}

	return nil
}

func ValidateLogin(login string) error {
	if !emailRegex.MatchString(login) && !phoneNumberRegex.MatchString(login) {
		return errors.NewErrorWithMessage("login is not valid")
	}

	return nil
}
