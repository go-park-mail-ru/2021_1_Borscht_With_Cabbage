package validation

import "github.com/borscht/backend/internal/models"

func ValidateUserRegistration(newUser models.User) error {
	err := ValidateEmail(newUser.Email)
	if err != nil {
		return err
	}
	err = ValidatePhoneNumber(newUser.Phone)
	if err != nil {
		return err
	}
	err = ValidateName(newUser.Name)
	if err != nil {
		return err
	}
	err = ValidatePassword(newUser.Password)
	if err != nil {
		return err
	}

	return nil
}

func ValidateUserLogin(user models.UserAuth) error {
	err := ValidateLogin(user.Login)
	if err != nil {
		return err
	}
	err = ValidatePassword(user.Password)
	if err != nil {
		return err
	}

	return nil
}

func ValidateUserEdit(user models.UserData) error {
	err := ValidateEmail(user.Email)
	if err != nil {
		return err
	}
	err = ValidatePhoneNumber(user.Phone)
	if err != nil {
		return err
	}
	err = ValidateName(user.Name)
	if err != nil {
		return err
	}
	err = ValidatePassword(user.Password)
	if err != nil {
		return err
	}

	return nil
}
