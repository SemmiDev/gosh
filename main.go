package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cache"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/websocket/v2"
	"github.com/jackc/pgx/v4/pgxpool"
)

type Program struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type ProgramDataStore struct {
	pool *pgxpool.Pool
}

func (p *ProgramDataStore) CreateProgram(ctx context.Context, program Program) (Program, error) {
	var id int64
	err := p.pool.QueryRow(ctx, "INSERT INTO program (name, description) VALUES ($1, $2) RETURNING id", program.Name, program.Description).Scan(&id)
	if err != nil {
		return Program{}, err
	}

	program.ID = id
	return program, nil
}

func (p *ProgramDataStore) GetProgramDetails(ctx context.Context, id int64) (Program, error) {
	var program Program
	err := p.pool.QueryRow(ctx, "SELECT id, name, description FROM program WHERE id = $1", id).Scan(&program.ID, &program.Name, &program.Description)
	return program, err
}

func (p *ProgramDataStore) SearchTerm(ctx context.Context, q string) ([]Program, error) {
	// english, indonesia, etc..
	rows, err := p.pool.Query(ctx, "SELECT id, name, description FROM program WHERE ts @@ to_tsquery('english', $1)", q)
	if err != nil {
		return nil, err
	}

	var programs []Program
	for rows.Next() {
		var program Program
		err = rows.Scan(&program.ID, &program.Name, &program.Description)
		fmt.Println("scan program: ", program)
		if err != nil {
			return nil, err
		}
		programs = append(programs, program)
	}

	return programs, nil
}

type HttpHandler struct {
	programDataStore *ProgramDataStore
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

	// change the space to underscore
	// because we use ts-vector
	q = strings.ReplaceAll(q, " ", "_")

	programs, err := h.programDataStore.SearchTerm(c.Context(), q)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Unable to search program",
		})
	}

	response := fiber.Map{
		"data": programs,
	}

	return c.Status(fiber.StatusOK).JSON(response)
}

func main() {
	pool, err := pgxpool.Connect(context.Background(), "postgres://root:secret@localhost/gosh")
	if err != nil {
		log.Fatalln("Unable to connect to database:", err)
	}
	defer pool.Close()

	err = pool.Ping(context.Background())
	if err != nil {
		log.Fatalln("Unable ping to database:", err)
	}

	programDataStore := &ProgramDataStore{pool: pool}
	httpHandler := &HttpHandler{programDataStore: programDataStore}

	app := fiber.New()

	app.Use(cors.New())
	app.Use(logger.New())

	app.Use("/ws", func(c *fiber.Ctx) error {
		// IsWebSocketUpgrade returns true if the client
		// requested upgrade to the WebSocket protocol.
		if websocket.IsWebSocketUpgrade(c) {
			c.Locals("allowed", true)
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})

	app.Get("/ws", websocket.New(func(c *websocket.Conn) {
		// c.Locals is added to the *websocket.Conn
		log.Println(c.Locals("allowed"))  // true
		log.Println(c.Params("id"))       // 123
		log.Println(c.Query("v"))         // 1.0
		log.Println(c.Cookies("session")) // ""

		// websocket.Conn bindings https://pkg.go.dev/github.com/fasthttp/websocket?tab=doc#pkg-index
		var (
			mt  int
			msg []byte
			err error
		)
		for {
			if mt, msg, err = c.ReadMessage(); err != nil {
				log.Println("read:", err)
				break
			}

			type searchRequest struct {
				Q string `json:"q"`
			}

			var req searchRequest
			err = json.Unmarshal(msg, &req)
			if err != nil {
				msg = []byte("Invalid request body")
				if err = c.WriteMessage(mt, msg); err != nil {
					log.Println("write:", err)
					break
				}
			}

			q := strings.ReplaceAll(strings.TrimSpace(req.Q), " ", "_")
			programs, err := programDataStore.SearchTerm(context.Background(), q)
			if err != nil {
				msg = []byte("Unable to searching program")
				if err = c.WriteMessage(mt, msg); err != nil {
					log.Println("write:", err)
					break
				}
			}

			response := fiber.Map{
				"data": []Program{},
			}

			if len(programs) != 0 {
				response["data"] = programs
			}

			msg, err = json.Marshal(response)
			if err != nil {
				msg = []byte("Unable to marshal response")
				if err = c.WriteMessage(mt, msg); err != nil {
					log.Println("write:", err)
					break
				}
			}

			if err = c.WriteMessage(mt, msg); err != nil {
				log.Println("write:", err)
				break
			}
		}

	}))

	app.Use(cache.New(cache.Config{
		Next: func(c *fiber.Ctx) bool {
			return strings.Contains(c.Route().Path, "/ws")
		},
	}))

	app.Get("/ping", func(c *fiber.Ctx) error {
		return c.SendString("pong 👋")
	})

	app.Post("/api/program", httpHandler.CreateProgramHandler)
	app.Get("/api/program", httpHandler.SearchProgram)
	app.Get("/api/program/:id", httpHandler.GetProgramDetailsHandler)

	err = app.Listen(":3000")
	if err != nil {
		log.Fatalln("Unable to start server:", err)
	}
}
