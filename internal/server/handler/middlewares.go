package handler

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
)

var (
	jwtSecret = ""
)

func AuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		tokenStr := c.Request().Header.Get("Authorization")
		if tokenStr == "" {
			resp := ErrorResponse{Errors: "missing token"}
			return echo.NewHTTPError(http.StatusUnauthorized, resp)
		}

		parts := strings.Split(tokenStr, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			resp := ErrorResponse{Errors: "invalid token format"}
			return echo.NewHTTPError(http.StatusUnauthorized, resp)
		}

		token, err := jwt.Parse(parts[1], func(token *jwt.Token) (any, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(jwtSecret), nil
		})

		if err != nil || !token.Valid {
			resp := ErrorResponse{Errors: "invalid token"}
			return echo.NewHTTPError(http.StatusUnauthorized, resp)
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			resp := ErrorResponse{Errors: "invalid claims"}
			return echo.NewHTTPError(http.StatusUnauthorized, resp)
		}

		userID, ok := claims["user_id"].(float64)
		if !ok {
			resp := ErrorResponse{Errors: "invalid user_id"}
			return echo.NewHTTPError(http.StatusUnauthorized, resp)
		}

		c.Set("user_id", uint(userID))
		return next(c)
	}
}
