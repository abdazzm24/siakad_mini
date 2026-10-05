// Package response memusatkan bentuk seluruh response JSON agar selalu seragam.
package response

import "github.com/gofiber/fiber/v2"

type Meta struct {
	CurrentPage int `json:"current_page"`
	PerPage     int `json:"per_page"`
	Total       int `json:"total"`
	LastPage    int `json:"last_page"`
}

type Body struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
	Meta    *Meta  `json:"meta,omitempty"`
	Errors  any    `json:"errors,omitempty"`
}

// LoginBody mengikuti format response login pada soal: token berada di level atas.
type LoginBody struct {
	Success     bool   `json:"success"`
	Message     string `json:"message"`
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
	Data        any    `json:"data"`
}

func Success(c *fiber.Ctx, message string, data any) error {
	return c.Status(fiber.StatusOK).JSON(Body{Success: true, Message: message, Data: data})
}

func Created(c *fiber.Ctx, message string, data any, location string) error {
	if location != "" {
		c.Set(fiber.HeaderLocation, location)
	}
	return c.Status(fiber.StatusCreated).JSON(Body{Success: true, Message: message, Data: data})
}

func Paginated(c *fiber.Ctx, message string, data any, meta Meta) error {
	return c.Status(fiber.StatusOK).JSON(Body{Success: true, Message: message, Data: data, Meta: &meta})
}

func NoContent(c *fiber.Ctx) error {
	return c.SendStatus(fiber.StatusNoContent)
}

// Error dipakai ErrorHandler. errs boleh nil.
func Error(c *fiber.Ctx, status int, message string, errs map[string][]string) error {
	body := Body{Success: false, Message: message}
	if len(errs) > 0 {
		body.Errors = errs
	}
	return c.Status(status).JSON(body)
}

// NewMeta menghitung last_page (minimal 1).
func NewMeta(page, perPage, total int) Meta {
	last := 1
	if perPage > 0 && total > 0 {
		last = (total + perPage - 1) / perPage
	}
	return Meta{CurrentPage: page, PerPage: perPage, Total: total, LastPage: last}
}