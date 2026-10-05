package handler

import (
	"github.com/gofiber/fiber/v2"

	"siakad-mini/app/request"
	"siakad-mini/app/response"
	"siakad-mini/app/service"
)

type AuthHandler struct{ svc *service.AuthService }

func NewAuthHandler(svc *service.AuthService) *AuthHandler { return &AuthHandler{svc: svc} }

// Login: POST /api/v1/auth/login
func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req request.LoginRequest

	if err := bindAndValidate(c, &req); err != nil {
		return err
	}
	req.Normalize()
	if err := validateStruct(&req); err != nil {
		return err
	}

	ctx, cancel := requestContext(c)
	defer cancel()

	result, err := h.svc.Login(ctx, req)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(response.LoginBody{
		Success:     true,
		Message:     "Login berhasil",
		AccessToken: result.AccessToken,
		TokenType:   result.TokenType,
		ExpiresIn:   result.ExpiresIn,
		Data:        result.User,
	})
}

// Me: GET /api/v1/auth/me
func (h *AuthHandler) Me(c *fiber.Ctx) error {
	au, err := authUser(c)
	if err != nil {
		return err
	}

	ctx, cancel := requestContext(c)
	defer cancel()

	me, err := h.svc.Me(ctx, au)
	if err != nil {
		return err
	}
	return response.Success(c, "Profil berhasil diambil", me)
}