package service

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"
	"sahamrakyat_test/database"
	"sahamrakyat_test/helpers"

	"github.com/go-playground/validator"
	"github.com/go-redis/cache/v8"
	"github.com/labstack/echo/v4"
)

func Create(c echo.Context) (*database.Orders, *echo.HTTPError) {
	logger := helpers.InitLogger()

	logger.Info(helpers.ApacheFormatLogger(c.Request().Method, c.Request().URL.Host, c.Request().Host, c.RealIP(), c.Request().UserAgent(), time.Now().String()))

	order := &database.Orders{}

	if err := c.Bind(order); err != nil {
		fmt.Println(err)
		return nil, echo.NewHTTPError(http.StatusInternalServerError, "Failed to bind order.")
	}

	if err := validator.New().Struct(order); err != nil {
		fmt.Println(err)

		validationErrors := []echo.Map{}

		if _, ok := err.(*validator.InvalidValidationError); ok {
			return nil, echo.NewHTTPError(http.StatusInternalServerError, "Invalid argument passed.")
		}

		for _, err := range err.(validator.ValidationErrors) {

			validationErrors = append(validationErrors, echo.Map{
				"namespace":       err.Namespace(),
				"field":           err.Field(),
				"structNamespace": err.StructNamespace(),
				"structField":     err.StructField(),
				"tag":             err.Tag(),
				"actualTag":       err.ActualTag(),
				"kind":            err.Kind(),
				"type":            err.Type(),
				"value":           err.Value(),
				"param":           err.Param(),
			})
		}

		return nil, echo.NewHTTPError(http.StatusBadRequest, echo.Map{
			"message": "Failed to validate order.",
			"error":   validationErrors,
		})
	}

	db := helpers.ConnectDatabase(os.Getenv("DB_HOST"), os.Getenv("DB_USER"), os.Getenv("DB_PASS"), os.Getenv("DB_NAME"))

	db.Create(order)

	logger.Info(fmt.Sprintf("[Response] %s %s %s %s %s", c.Response().Header().Get("Method"), c.Response().Header().Get("Host"), c.Response().Header().Get("RemoteAddr"), c.Response().Header().Get("UserAgent"), c.Response().Header().Get("Time")))

	return order, nil
}

func GetAll(c echo.Context) (*[]database.Orders, *echo.HTTPError) {
	logger := helpers.InitLogger()

	logger.Info(helpers.ApacheFormatLogger(c.Request().Method, c.Request().URL.Host, c.Request().Host, c.RealIP(), c.Request().UserAgent(), time.Now().String()))

	var take int
	var skip int

	takeQuery := c.QueryParam("take")

	if takeQuery != "" {
		if val, err := strconv.ParseInt(takeQuery, 10, 32); err != nil {
			return nil, echo.NewHTTPError(http.StatusBadRequest, "Failed to parse take.")
		} else {
			take = int(val)
		}
	} else {
		take = -1
	}

	skipQuery := c.QueryParam("skip")

	if skipQuery != "" {
		if val, err := strconv.ParseInt(skipQuery, 10, 32); err != nil {
			return nil, echo.NewHTTPError(http.StatusBadRequest, "Failed to parse skip.")
		} else {
			skip = int(val)
		}
	} else {
		skip = -1
	}

	orders := &[]database.Orders{}

	db := helpers.ConnectDatabase(os.Getenv("DB_HOST"), os.Getenv("DB_USER"), os.Getenv("DB_PASS"), os.Getenv("DB_NAME"))

	db.Limit(take).Offset(skip).Find(orders)

	cacheClient := helpers.InitRedisCache()

	if err := cacheClient.Get(c.Request().Context(), "orders", orders); err == nil {
		return orders, nil
	} else if err != nil {
		if err := cacheClient.Set(&cache.Item{
			Ctx:   c.Request().Context(),
			Key:   "orders",
			Value: orders,
		}); err != nil {
			fmt.Printf("Failed to cache orders: %v", err)
		}
	}

	logger.Info(fmt.Sprintf("[Response] %s %s %s %s %s", c.Response().Header().Get("Method"), c.Response().Header().Get("Host"), c.Response().Header().Get("RemoteAddr"), c.Response().Header().Get("UserAgent"), c.Response().Header().Get("Time")))

	return orders, nil
}

func Get(c echo.Context) (*database.Orders, *echo.HTTPError) {
	logger := helpers.InitLogger()

	logger.Info(helpers.ApacheFormatLogger(c.Request().Method, c.Request().URL.Host, c.Request().Host, c.RealIP(), c.Request().UserAgent(), time.Now().String()))

	if c.Param("id") == "" {
		return nil, echo.NewHTTPError(http.StatusBadRequest, "Order id is required.")
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)

	if err != nil {
		return nil, echo.NewHTTPError(http.StatusBadRequest, "Failed to parse order id.")
	}

	order := &database.Orders{
		ID: uint(id),
	}

	if err := c.Bind(order); err != nil {
		return nil, echo.NewHTTPError(http.StatusInternalServerError, "Failed to bind order.")
	}

	db := helpers.ConnectDatabase(os.Getenv("DB_HOST"), os.Getenv("DB_USER"), os.Getenv("DB_PASS"), os.Getenv("DB_NAME"))

	if result := db.Find(order); result.Error != nil || result.RowsAffected < 0 {
		return nil, echo.NewHTTPError(http.StatusNotFound, "Order not found.")
	}

	cacheClient := helpers.InitRedisCache()

	if err := cacheClient.Get(c.Request().Context(), fmt.Sprintf("order:%d", order.ID), order); err == nil {
		return order, nil
	} else if err != nil {
		if err := cacheClient.Set(&cache.Item{
			Ctx:   c.Request().Context(),
			Key:   fmt.Sprintf("order:%d", order.ID),
			Value: order,
		}); err != nil {
			fmt.Printf("Failed to cache order: %v", err)
		}
	}

	logger.Info(fmt.Sprintf("[Response] %s %s %s %s %s", c.Response().Header().Get("Method"), c.Response().Header().Get("Host"), c.Response().Header().Get("RemoteAddr"), c.Response().Header().Get("UserAgent"), c.Response().Header().Get("Time")))

	return order, nil
}

func Update(c echo.Context) (*database.Orders, *echo.HTTPError) {
	logger := helpers.InitLogger()

	logger.Info(helpers.ApacheFormatLogger(c.Request().Method, c.Request().URL.Host, c.Request().Host, c.RealIP(), c.Request().UserAgent(), time.Now().String()))

	if c.Param("id") == "" {
		return nil, echo.NewHTTPError(http.StatusBadRequest, "Order id is required.")
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)

	if err != nil {
		return nil, echo.NewHTTPError(http.StatusBadRequest, "Failed to parse order id.")
	}

	order := &database.Orders{
		ID: uint(id),
	}

	if err := c.Bind(order); err != nil {
		fmt.Println(err)
		return nil, echo.NewHTTPError(http.StatusInternalServerError, "Failed to bind order.")
	}

	db := helpers.ConnectDatabase(os.Getenv("DB_HOST"), os.Getenv("DB_USER"), os.Getenv("DB_PASS"), os.Getenv("DB_NAME"))

	if result := db.Find(order); result.Error != nil || result.RowsAffected < 0 {
		return nil, echo.NewHTTPError(http.StatusNotFound, "Order not found.")
	}

	cacheClient := helpers.InitRedisCache()

	if err := cacheClient.Delete(c.Request().Context(), fmt.Sprintf("order:%d", order.ID)); err == nil {
		cacheClient.Delete(c.Request().Context(), fmt.Sprintf("order:%d", id))
	}

	db.Save(order)

	logger.Info(fmt.Sprintf("[Response] %s %s %s %s %s", c.Response().Header().Get("Method"), c.Response().Header().Get("Host"), c.Response().Header().Get("RemoteAddr"), c.Response().Header().Get("UserAgent"), c.Response().Header().Get("Time")))

	return order, nil
}

func Delete(c echo.Context) (*database.Orders, *echo.HTTPError) {
	logger := helpers.InitLogger()

	logger.Info(helpers.ApacheFormatLogger(c.Request().Method, c.Request().URL.Host, c.Request().Host, c.RealIP(), c.Request().UserAgent(), time.Now().String()))

	if c.Param("id") == "" {
		return nil, echo.NewHTTPError(http.StatusBadRequest, "Order id is required.")
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)

	if err != nil {
		return nil, echo.NewHTTPError(http.StatusBadRequest, "Failed to parse order id.")
	}

	order := &database.Orders{
		ID: uint(id),
	}

	if err := c.Bind(order); err != nil {
		fmt.Println(err)
		return nil, echo.NewHTTPError(http.StatusInternalServerError, "Failed to bind order.")
	}

	db := helpers.ConnectDatabase(os.Getenv("DB_HOST"), os.Getenv("DB_USER"), os.Getenv("DB_PASS"), os.Getenv("DB_NAME"))

	if result := db.Find(order); result.Error != nil || result.RowsAffected < 0 {
		return nil, echo.NewHTTPError(http.StatusNotFound, "Order not found.")
	}

	cacheClient := helpers.InitRedisCache()

	if err := cacheClient.Delete(c.Request().Context(), fmt.Sprintf("order:%d", order.ID)); err == nil {
		cacheClient.Delete(c.Request().Context(), fmt.Sprintf("order:%d", id))
	}

	db.Delete(order)

	logger.Info(fmt.Sprintf("[Response] %s %s %s %s %s", c.Response().Header().Get("Method"), c.Response().Header().Get("Host"), c.Response().Header().Get("RemoteAddr"), c.Response().Header().Get("UserAgent"), c.Response().Header().Get("Time")))

	return order, nil
}
