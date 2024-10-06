package middlewares

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"message-automation/src/models/base"
	"net/http"
)

func RecoverMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		defer func() {
			if r := recover(); r != nil {
				switch e := r.(type) {
				case *base.NotFoundError:
					base.Log(fmt.Sprintf("NotFoundError: %s", e.Error()))
					_ = c.JSON(http.StatusNotFound, base.Response[base.Error]{
						Success: false,
						Data: base.Error{
							Code:    404,
							Message: e.Error(),
						},
					})
				case *base.BadRequestError:
					base.Log(fmt.Sprintf("BadRequestError: %s", e.Error()))
					_ = c.JSON(http.StatusBadRequest, base.Response[base.Error]{
						Success: false,
						Data: base.Error{
							Code:    400,
							Message: e.Error(),
						},
					})
				default:
					base.Log(fmt.Sprintf("InternalServerError: %s", e))
					_ = c.JSON(http.StatusInternalServerError, base.Response[base.Error]{
						Success: false,
						Data: base.Error{
							Code:    500,
							Message: fmt.Sprintf("%s", e),
						},
					})
				}
			}
		}()
		return next(c)
	}
}
