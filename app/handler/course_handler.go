package handler

import (
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"

	"siakad-mini/app/apperror"
	"siakad-mini/app/model"
	"siakad-mini/app/response"
	"siakad-mini/app/service"
)

type CourseHandler struct{ svc *service.CourseService }

func NewCourseHandler(svc *service.CourseService) *CourseHandler { return &CourseHandler{svc: svc} }

// List: GET /api/v1/courses?semester=&search=&available=true
func (h *CourseHandler) List(c *fiber.Ctx) error {
	filter := model.CourseFilter{
		Search:        strings.TrimSpace(c.Query("search")),
		AvailableOnly: isTrue(c.Query("available")),
	}
	if raw := strings.TrimSpace(c.Query("semester")); raw != "" {
		semester, err := strconv.Atoi(raw)
		if err != nil || semester < 1 {
			return apperror.FieldError("semester", "Semester harus berupa angka positif")
		}
		filter.Semester = &semester
	}

	ctx, cancel := requestContext(c)
	defer cancel()

	courses, err := h.svc.List(ctx, filter)
	if err != nil {
		return err
	}
	return response.Success(c, "Data mata kuliah berhasil diambil", courses)
}

func isTrue(value string) bool {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "true", "1", "yes":
		return true
	}
	return false
}