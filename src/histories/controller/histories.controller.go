package controller

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"sahamrakyat_test/histories/service"
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
		"message":    "Successfully created new history.",
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
		"message":    "Successfully get all histories.",
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
		"message":    "Successfully get history.",
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
		"message":    "Successfully updated history.",
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
		"message":    "Successfully deleted history.",
		"data":       data,
	})
}
