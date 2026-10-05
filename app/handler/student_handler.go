package handler

import (
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"

	"siakad-mini/app/model"
	"siakad-mini/app/request"
	"siakad-mini/app/response"
	"siakad-mini/app/service"
)

const (
	defaultPerPage = 10
	maxPerPage     = 50
)

type StudentHandler struct{ svc *service.StudentService }

func NewStudentHandler(svc *service.StudentService) *StudentHandler {
	return &StudentHandler{svc: svc}
}

// parseListQuery membaca query string dan memberi nilai bawaan yang aman.
func parseListQuery(c *fiber.Ctx) model.StudentListQuery {
	q := model.StudentListQuery{
		Page:    c.QueryInt("page", 1),
		PerPage: c.QueryInt("per_page", defaultPerPage),
		Prodi:   strings.TrimSpace(c.Query("prodi")),
		Search:  strings.TrimSpace(c.Query("search")),
		Sort:    strings.TrimSpace(c.Query("sort")),
	}
	if q.Page < 1 {
		q.Page = 1
	}
	if q.PerPage < 1 {
		q.PerPage = defaultPerPage
	}
	if q.PerPage > maxPerPage {
		q.PerPage = maxPerPage
	}
	if raw := strings.TrimSpace(c.Query("angkatan")); raw != "" {
		if angkatan, err := strconv.Atoi(raw); err == nil {
			q.Angkatan = &angkatan
		}
	}
	return q
}

// List: GET /api/v1/students
func (h *StudentHandler) List(c *fiber.Ctx) error {
	q := parseListQuery(c)

	ctx, cancel := requestContext(c)
	defer cancel()

	students, total, err := h.svc.List(ctx, q)
	if err != nil {
		return err
	}
	return response.Paginated(c, "Data mahasiswa berhasil diambil", students,
		response.NewMeta(q.Page, q.PerPage, total))
}

// Create: POST /api/v1/students
func (h *StudentHandler) Create(c *fiber.Ctx) error {
	var req request.CreateStudentRequest
	if err := bindAndValidate(c, &req); err != nil {
		return err
	}
	req.Normalize()
	if err := validateStruct(&req); err != nil {
		return err
	}

	ctx, cancel := requestContext(c)
	defer cancel()

	student, err := h.svc.Create(ctx, req)
	if err != nil {
		return err
	}
	return response.Created(c, "Mahasiswa berhasil ditambahkan", student,
		"/api/v1/students/"+strconv.FormatInt(student.ID, 10))
}

// Show: GET /api/v1/students/{id}[?tahun_akademik=2026/2027-Ganjil]
func (h *StudentHandler) Show(c *fiber.Ctx) error {
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

	detail, err := h.svc.GetDetail(ctx, au, id, strings.TrimSpace(c.Query("tahun_akademik")))
	if err != nil {
		return err
	}
	return response.Success(c, "Detail mahasiswa berhasil diambil", detail)
}

// Update: PUT /api/v1/students/{id}
func (h *StudentHandler) Update(c *fiber.Ctx) error {
	id, err := pathID(c)
	if err != nil {
		return err
	}

	var req request.UpdateStudentRequest
	if err := bindAndValidate(c, &req); err != nil {
		return err
	}
	req.Normalize()
	if err := validateStruct(&req); err != nil {
		return err
	}
	if req.IsEmpty() {
		return emptyUpdateError()
	}

	ctx, cancel := requestContext(c)
	defer cancel()

	student, err := h.svc.Update(ctx, id, req)
	if err != nil {
		return err
	}
	return response.Success(c, "Mahasiswa berhasil diperbarui", student)
}

// Delete: DELETE /api/v1/students/{id} (soft delete)
func (h *StudentHandler) Delete(c *fiber.Ctx) error {
	id, err := pathID(c)
	if err != nil {
		return err
	}

	ctx, cancel := requestContext(c)
	defer cancel()

	if err := h.svc.Delete(ctx, id); err != nil {
		return err
	}
	return response.NoContent(c)
}