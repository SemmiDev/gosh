package program

import (
	"github.com/gofiber/fiber/v2"
)

type HttpHandler struct {
	programDataStore *ProgramDataStore
	searchService    *SearchService
}

func NewHttpHandler(
	programDataStore *ProgramDataStore,
	searchService *SearchService,
) *HttpHandler {
	return &HttpHandler{
		programDataStore: programDataStore,
		searchService:    searchService,
	}
}

type CreateProgramRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

func (h *HttpHandler) CreateProgramHandler(c *fiber.Ctx) error {
	var req CreateProgramRequest
	err := c.BodyParser(&req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
		})
	}

	program := Program{
		Name:        req.Name,
		Description: req.Description,
	}

	program, err = h.programDataStore.CreateProgram(c.Context(), program)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Unable to create program",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Program created successfully",
		"data":    program,
	})
}

func (h *HttpHandler) GetProgramDetailsHandler(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid program id",
		})
	}

	program, err := h.programDataStore.GetProgramDetails(c.Context(), int64(id))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Program not found",
		})
	}

	response := fiber.Map{
		"data": program,
	}

	return c.Status(fiber.StatusOK).JSON(response)
}

func (h *HttpHandler) SearchProgram(c *fiber.Ctx) error {
	q := c.Query("q")
	if q == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid query",
		})
	}

	programs := h.searchService.Search(c.Context(), q)
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": programs,
	})
}
