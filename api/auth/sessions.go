package auth

import (
	"backend/api/domain"
	"backend/api/domain/user"
	errors "backend/models"
)


// будет использоваться для проверки уникальности сессии при создании и для проверки авторизации на сайте в целом
func CheckSession(sessionToCheck string, context *domain.CustomContext) (string, bool) {
	number, isItExists := (*context.Sessions)[sessionToCheck]
	if !isItExists {
		return "", false
	}
	return number, true
}

func GetUser(context *domain.CustomContext) (user.User, error) {
	sessionError := errors.Authorization("not authorization")
	sessionError.Description = "session error"
	session, err := context.Cookie(domain.SessionCookie)

	if err != nil {
		return user.User{}, sessionError
	}

	phone, ok := CheckSession(session.Value, context)
	if !ok {
		return user.User{}, sessionError
	}

	for _, user := range *context.Users {
		if user.Phone == phone {
			return user, nil
		}
	}

	return user.User{}, errors.Authorization("user not found")
}
