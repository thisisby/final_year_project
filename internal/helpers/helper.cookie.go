package helpers

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

func WriteCookie(c echo.Context, name string, value string) {
	cookie := new(http.Cookie)
	cookie.Name = name
	cookie.Value = value
	cookie.Path = "/"
	cookie.HttpOnly = true
	c.SetCookie(cookie)
}

func ReadCookie(c echo.Context, name string) (string, error) {
	cookie, err := c.Cookie(name)
	if err != nil {
		return "", err
	}
	return cookie.Value, nil
}
