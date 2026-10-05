package handler

import (
	"strconv"

	"github.com/gofiber/fiber/v2"

	"siakad-mini/app/request"
	"siakad-mini/app/response"
	"siakad-mini/app/service"
)

type EnrollmentHandler struct{ svc *service.EnrollmentService }

func NewEnrollmentHandler(svc *service.EnrollmentService) *EnrollmentHandler {
	return &EnrollmentHandler{svc: svc}
}

// Create: POST /api/v1/enrollments
func (h *EnrollmentHandler) Create(c *fiber.Ctx) error {
	au, err := authUser(c)
	if err != nil {
		return err
	}

	var req request.CreateEnrollmentRequest
	if err := bindAndValidate(c, &req); err != nil {
		return err
	}
	req.Normalize()
	if err := validateStruct(&req); err != nil {
		return err
	}

	ctx, cancel := requestContext(c)
	defer cancel()

	result, err := h.svc.Create(ctx, au, req)
	if err != nil {
		return err
	}
	return response.Created(c, "Mata kuliah berhasil ditambahkan ke KRS", result,
		"/api/v1/enrollments/"+strconv.FormatInt(result.ID, 10))
}

// Delete: DELETE /api/v1/enrollments/{id}
func (h *EnrollmentHandler) Delete(c *fiber.Ctx) error {
	au, err := authUser(c)
	if err != nil {
		return err
	}
	id, err := pathID(c)
	if err != nil {
		return err
	}

	ctx, cancel := requestContext(c)
	defer cancel()

	if err := h.svc.Delete(ctx, au, id); err != nil {
		return err
	}
	return response.NoContent(c)
}