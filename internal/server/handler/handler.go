package handler

import (
	"net/http"
	"os"
	"time"

	"github.com/0x0FACED/merch-shop/config"
	"github.com/0x0FACED/merch-shop/internal/model"
	"github.com/0x0FACED/merch-shop/internal/server/validator"
	"github.com/0x0FACED/merch-shop/internal/service"
	"github.com/0x0FACED/merch-shop/pkg/logger"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	userService *service.MerchService

	logger *logger.ZapLogger
	config *config.ServerConfig
}

func NewHandler(u *service.MerchService, l *logger.ZapLogger, cfg *config.ServerConfig) *Handler {
	jwtSecret = os.Getenv("JWT_SECRET_KEY")

	return &Handler{
		userService: u,
		logger:      l,
		config:      cfg,
	}
}

func (h *Handler) SetupRoutes(e *echo.Echo) {
	e.POST("/api/auth", h.AuthUser) // Аутентификация юзера
	group := e.Group("/api", AuthMiddleware)

	group.GET("/info", h.GetUserInfo)   // Получаем всю инфу о юзере (транзакции, баланс, инвентарь)
	group.GET("/buy/:item", h.BuyItem)  // Делаем покупку предмета юзером (why GET?)
	group.POST("/sendCoin", h.SendCoin) // отправка монет кому-либо
}

func (h *Handler) AuthUser(c echo.Context) error {
	var req AuthRequest

	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request")
	}

	if err := c.Validate(&req); err != nil {
		if validationErrs, ok := err.(*validator.ValidationErrorsResponse); ok {
			return c.JSON(http.StatusBadRequest, validationErrs)
		}
		resp := ErrorResponse{Errors: err.Error()}
		return echo.NewHTTPError(http.StatusBadRequest, resp)

	}

	params := model.AuthUserParams{
		Username: req.Username,
		Password: req.Password,
	}

	ctx := c.Request().Context()

	user, err := h.userService.AuthUser(ctx, params)
	if err != nil {
		resp := ErrorResponse{Errors: err.Error()}
		return echo.NewHTTPError(MapServiceErrorToStatusCode(err), resp)
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(time.Hour * 72).Unix(),
	})

	tokenString, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		resp := ErrorResponse{Errors: "incorrect login or password"}
		return echo.NewHTTPError(MapServiceErrorToStatusCode(err), resp)
	}

	resp := AuthResponse{
		Token: tokenString,
	}

	return c.JSON(http.StatusOK, resp)
}

func (h *Handler) GetUserInfo(c echo.Context) error {
	userID := c.Get("user_id").(uint)

	ctx := c.Request().Context()

	params := model.GetUserInfoParams{
		ID: userID,
	}

	userInfo, err := h.userService.GetUserInfo(ctx, params)
	if err != nil {
		return echo.NewHTTPError(MapServiceErrorToStatusCode(err), "failed to get user info")
	}

	resp := InfoResponse{
		Coins:       userInfo.Coins,
		Inventory:   userInfo.Inventory,
		CoinHistory: userInfo.CoinHistory,
	}

	return c.JSON(http.StatusOK, resp)
}

func (h *Handler) BuyItem(c echo.Context) error {
	userID := c.Get("user_id").(uint)
	item := c.Param("item")

	ctx := c.Request().Context()

	params := model.BuyItemParams{
		UserID: userID,
		Item:   item,
	}
	if err := h.userService.BuyItem(ctx, params); err != nil {
		resp := ErrorResponse{Errors: err.Error()}
		return echo.NewHTTPError(MapServiceErrorToStatusCode(err), resp)
	}

	return c.NoContent(http.StatusOK)
}

func (h *Handler) SendCoin(c echo.Context) error {
	userID := c.Get("user_id").(uint)

	var req SendCoinRequest

	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request")
	}

	if err := c.Validate(&req); err != nil {
		if validationErrs, ok := err.(*validator.ValidationErrorsResponse); ok {
			return c.JSON(http.StatusBadRequest, validationErrs)
		}
		resp := ErrorResponse{Errors: err.Error()}
		return echo.NewHTTPError(http.StatusBadRequest, resp)

	}

	params := model.SendCoinParams{
		FromUser: userID,
		ToUser:   req.ToUser,
		Amount:   req.Amount,
	}

	ctx := c.Request().Context()

	if err := h.userService.SendCoin(ctx, params); err != nil {
		resp := ErrorResponse{Errors: err.Error()}
		return echo.NewHTTPError(MapServiceErrorToStatusCode(err), resp)
	}

	return c.NoContent(http.StatusOK)
}
