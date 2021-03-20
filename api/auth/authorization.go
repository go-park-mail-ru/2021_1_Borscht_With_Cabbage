package auth

import (
	"backend/api/domain"
	errors "backend/models"
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
)

func LogoutUser(c echo.Context) error {
	cc := c.(*domain.CustomContext)
	cook, err := cc.Cookie(domain.SessionCookie)
	if err != nil {
		sendErr := errors.Create(http.StatusUnauthorized, "error with request data")
		return cc.SendERR(sendErr)
	}

	fmt.Println(*cc.Sessions)
	session := cook.Value

	_, ok := (*cc.Sessions)[session]
	if ok {
		delete(*cc.Sessions, session)
	}

	fmt.Println(*cc.Sessions)

	DeleteResponseCookie(c)

	return cc.SendOK("")
}
