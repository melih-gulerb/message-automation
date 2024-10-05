package middlewares

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"log"
	"message-automation/src/models/base"
	"net/http"
)

func RecoverMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		defer func() {
			if r := recover(); r != nil {
				switch e := r.(type) {
				case *base.NotFoundError:
					log.Printf("NotFoundError: %s", e.Error())
					_ = c.JSON(http.StatusNotFound, map[string]interface{}{
						"error":   "Not Found",
						"message": e.Error(),
					})
				case *base.BadRequestError:
					log.Printf("BadRequestError: %s", e.Error())
					_ = c.JSON(http.StatusBadRequest, map[string]interface{}{
						"error":   "Bad Request",
						"message": e.Error(),
					})
				default:
					log.Printf("InternalServerError: %s", e)
					_ = c.JSON(http.StatusInternalServerError, map[string]interface{}{
						"error":   "Internal Server Error",
						"message": fmt.Sprintf("exception: %v", e),
					})
				}
			}
		}()
		return next(c)
	}
}
