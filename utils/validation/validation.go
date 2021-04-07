package validation

import (
	"github.com/borscht/backend/utils/errors"
	"regexp"
)

var emailRegex = regexp.MustCompile("^([A-Za-z0-9_\\-.])+@([A-Za-z0-9_\\-.])+\\.([A-Za-z]{2,4})$")
var phoneNumberRegex = regexp.MustCompile("^[0-9]{11}$")

func ValidateEmail(email string) error {
	if len(email) == 0 {
		sendErr := errors.BadRequestError("email can't be empty")
		return sendErr
	}

	if !emailRegex.MatchString(email) {
		sendErr := errors.BadRequestError("email is not valid")
		return sendErr
	}

	return nil
}

func ValidatePhoneNumber(number string) error {
	if len(number) == 0 {
		sendErr := errors.BadRequestError("phone number can't be empty")
		return sendErr
	}

	if !phoneNumberRegex.MatchString(number) {
		sendErr := errors.BadRequestError("phone number is not valid")
		return sendErr
	}

	return nil
}

func ValidatePassword(password string) error {
	passwordLength := len(password)

	if passwordLength == 0 {
		sendErr := errors.BadRequestError("password can't be empty")
		return sendErr
	}

	if passwordLength > 30 || passwordLength < 6 {
		sendErr := errors.BadRequestError("password must be 6-30 symbols")
		return sendErr
	}

	return nil
}

func ValidateName(name string) error {
	if len(name) == 0 {
		sendErr := errors.BadRequestError("name can't be empty")
		return sendErr
	}

	return nil
}

func ValidateLogin(login string) error {
	if !emailRegex.MatchString(login) && !phoneNumberRegex.MatchString(login) {
		sendErr := errors.BadRequestError("login is not valid")
		return sendErr
	}

	return nil
}
