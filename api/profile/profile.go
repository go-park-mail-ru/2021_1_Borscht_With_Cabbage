package profile

import (
	"backend/api"
	"backend/api/auth"
	"encoding/json"
	"github.com/labstack/echo/v4"
	"net/http"
)

func GetUserData(c echo.Context) error {
	cc := c.(*api.CustomContext)

	user, err := auth.GetUser(cc)
	if err != nil {
		return err
	}

	userData, err := json.Marshal(user)
	if err != nil {
		return err
	}
	return c.String(http.StatusOK, string(userData))
}
