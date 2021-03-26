package middleware

import (
	"net/http"

	"github.com/borscht/backend/config"
	"github.com/labstack/echo/v4/middleware"
)

var CORS = middleware.CORSWithConfig(middleware.CORSConfig{
	AllowOrigins:     []string{config.Client, "http://127.0.0.1:3000"},
	AllowMethods:     []string{http.MethodGet, http.MethodPut, http.MethodPost, http.MethodDelete},
	AllowCredentials: true,
})
