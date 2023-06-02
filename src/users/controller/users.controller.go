package controller

import (
	"net/http"
	"sahamrakyat_test/users/service"

	"github.com/labstack/echo/v4"
)

func Create(c echo.Context) error {
	data, err := service.Create(c)

	if err != nil {
		return c.JSON(err.Code, echo.Map{
			"statusCode": err.Code,
			"message":    err.Message,
		})
	}

	return c.JSON(http.StatusCreated, echo.Map{
		"statusCode": http.StatusCreated,
		"message":    "Successfully created new user.",
		"data":       data,
	})
}

func GetAll(c echo.Context) error {
	data, err := service.GetAll(c)

	if err != nil {
		return c.JSON(err.Code, echo.Map{
			"statusCode": err.Code,
			"message":    err.Message,
		})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"statusCode": http.StatusOK,
		"message":    "Successfully get all users.",
		"data":       data,
	})
}

func Get(c echo.Context) error {
	data, err := service.Get(c)

	if err != nil {
		return c.JSON(err.Code, echo.Map{
			"statusCode": err.Code,
			"message":    err.Message,
		})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"statusCode": http.StatusOK,
		"message":    "Successfully get user.",
		"data":       data,
	})
}

func Update(c echo.Context) error {
	data, err := service.Update(c)

	if err != nil {
		return c.JSON(err.Code, echo.Map{
			"statusCode": err.Code,
			"message":    err.Message,
		})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"statusCode": http.StatusOK,
		"message":    "Successfully updated user.",
		"data":       data,
	})
}

func Delete(c echo.Context) error {
	data, err := service.Delete(c)

	if err != nil {
		return c.JSON(err.Code, echo.Map{
			"statusCode": err.Code,
			"message":    err.Message,
		})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"statusCode": http.StatusOK,
		"message":    "Successfully deleted user.",
		"data":       data,
	})
}
