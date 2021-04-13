package validation

import (
	"github.com/borscht/backend/internal/models"
	"testing"
)

func TestValidateEmail(t *testing.T) {
	email := "dashamail.ru"
	err := ValidateEmail(email)
	if err == nil {
		t.Errorf("email validation error")
		return
	}
}

func TestValidateEmailEmpty(t *testing.T) {
	email := ""
	err := ValidateEmail(email)
	if err == nil {
		t.Errorf("email validation error")
		return
	}
}

func TestValidateLogin(t *testing.T) {
	login := "dashamail.ru"
	err := ValidateLogin(login)
	if err == nil {
		t.Errorf("login validation error")
		return
	}
}

func TestValidateLoginEmpty(t *testing.T) {
	login := ""
	err := ValidateLogin(login)
	if err == nil {
		t.Errorf("login validation error")
		return
	}
}

func TestValidateName(t *testing.T) {
	name := ""
	err := ValidateName(name)
	if err == nil {
		t.Errorf("name validation error")
		return
	}
}

func TestValidatePassword(t *testing.T) {
	password := "1111"
	err := ValidatePassword(password)
	if err == nil {
		t.Errorf("password validation error")
		return
	}
}

func TestValidatePasswordEmpty(t *testing.T) {
	password := ""
	err := ValidatePassword(password)
	if err == nil {
		t.Errorf("password validation error")
		return
	}
}

func TestValidatePhoneNumber(t *testing.T) {
	number := "1as1"
	err := ValidatePhoneNumber(number)
	if err == nil {
		t.Errorf("number validation error")
		return
	}
}

func TestValidatePhoneNumberEmpty(t *testing.T) {
	number := ""
	err := ValidatePhoneNumber(number)
	if err == nil {
		t.Errorf("number validation error")
		return
	}
}

func TestValidateSignIn_Login(t *testing.T) {
	err := ValidateSignIn("323", "11111111")
	if err == nil {
		t.Errorf("login validation error")
		return
	}
}

func TestValidateSignIn_Password(t *testing.T) {
	err := ValidateSignIn("dasha@mail.ru", "111")
	if err == nil {
		t.Errorf("password validation error")
		return
	}
}

func TestValidateUserRegistration(t *testing.T) {
	user := models.User{
		Email:    "dasha@mail.ru",
		Name:     "Daria",
		Phone:    "81111111111",
		Password: "1111111",
	}
	err := ValidateUserRegistration(user)
	if err != nil {
		t.Errorf("user registration validation error")
		return
	}
}

func TestValidateUserRegistration_Number(t *testing.T) {
	user := models.User{
		Email:    "dasha@mail.ru",
		Name:     "Daria",
		Phone:    "232",
		Password: "1111111",
	}
	err := ValidateUserRegistration(user)
	if err == nil {
		t.Errorf("number validation error")
		return
	}
}

func TestValidateUserRegistration_Email(t *testing.T) {
	user := models.User{
		Email:    "dashamailru",
		Name:     "Daria",
		Phone:    "89111111111",
		Password: "1111111",
	}
	err := ValidateUserRegistration(user)
	if err == nil {
		t.Errorf("email validation error")
		return
	}
}

func TestValidateUserRegistration_Name(t *testing.T) {
	user := models.User{
		Email:    "dasha@mail.ru",
		Name:     "",
		Phone:    "89111111111",
		Password: "1111111",
	}
	err := ValidateUserRegistration(user)
	if err == nil {
		t.Errorf("name validation error")
		return
	}
}

func TestValidateUserRegistration_Password(t *testing.T) {
	user := models.User{
		Email:    "dasha@mail.ru",
		Name:     "Daria",
		Phone:    "89111111111",
		Password: "111",
	}
	err := ValidateUserRegistration(user)
	if err == nil {
		t.Errorf("password validation error")
		return
	}
}

func TestValidateRestRegistration(t *testing.T) {
	restaurant := models.RestaurantInfo{
		AdminEmail:    "dasha@mail.ru",
		AdminPhone:    "89111111111",
		AdminPassword: "1111111",
		Title:         "rest1",
	}
	err := ValidateRestRegistration(restaurant)
	if err != nil {
		t.Errorf("restaurant registration validation error")
		return
	}
}

func TestValidateRestRegistration_Email(t *testing.T) {
	restaurant := models.RestaurantInfo{
		AdminEmail:    "dashamail.ru",
		AdminPhone:    "89111111111",
		AdminPassword: "1111111",
		Title:         "rest1",
	}
	err := ValidateRestRegistration(restaurant)
	if err == nil {
		t.Errorf("email validation error")
		return
	}
}

func TestValidateRestRegistration_Phone(t *testing.T) {
	restaurant := models.RestaurantInfo{
		AdminEmail:    "dasha@mail.ru",
		AdminPhone:    "891111",
		AdminPassword: "1111111",
		Title:         "rest1",
	}
	err := ValidateRestRegistration(restaurant)
	if err == nil {
		t.Errorf("email validation error")
		return
	}
}

func TestValidateRestRegistration_Password(t *testing.T) {
	restaurant := models.RestaurantInfo{
		AdminEmail:    "dasha@mail.ru",
		AdminPhone:    "89111111111",
		AdminPassword: "111",
		Title:         "rest1",
	}
	err := ValidateRestRegistration(restaurant)
	if err == nil {
		t.Errorf("password validation error")
		return
	}
}

func TestValidateRestRegistration_Title(t *testing.T) {
	restaurant := models.RestaurantInfo{
		AdminEmail:    "dasha@mail.ru",
		AdminPhone:    "89111111111",
		AdminPassword: "1111111",
		Title:         "",
	}
	err := ValidateRestRegistration(restaurant)
	if err == nil {
		t.Errorf("title validation error")
		return
	}
}

func TestValidateUserEdit(t *testing.T) {
	user := models.UserData{
		Email:    "dasha@mail.ru",
		Name:     "Daria",
		Phone:    "89111111111",
		Password: "111111111",
	}
	err := ValidateUserEdit(user)
	if err != nil {
		t.Errorf("user edit validation error")
		return
	}
}

func TestValidateUserEdit_Email(t *testing.T) {
	user := models.UserData{
		Email:    "dashamail.ru",
		Name:     "Daria",
		Phone:    "89111111111",
		Password: "11111111",
	}
	err := ValidateUserEdit(user)
	if err == nil {
		t.Errorf("user edit validation error")
		return
	}
}

func TestValidateUserEdit_Name(t *testing.T) {
	user := models.UserData{
		Email:    "dasha@mail.ru",
		Name:     "",
		Phone:    "89111111111",
		Password: "11111111",
	}
	err := ValidateUserEdit(user)
	if err == nil {
		t.Errorf("user edit validation error")
		return
	}
}

func TestValidateUserEdit_Phone(t *testing.T) {
	user := models.UserData{
		Email:    "dasha@mail.ru",
		Name:     "Daria",
		Phone:    "891111",
		Password: "11111111",
	}
	err := ValidateUserEdit(user)
	if err == nil {
		t.Errorf("user edit validation error")
		return
	}
}

func TestValidateUserEdit_Password(t *testing.T) {
	user := models.UserData{
		Email:    "dasha@mail.ru",
		Name:     "Daria",
		Phone:    "89111111111",
		Password: "111",
	}
	err := ValidateUserEdit(user)
	if err == nil {
		t.Errorf("user edit validation error")
		return
	}
}
