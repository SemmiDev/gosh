package internal

import (
	"context"
	"encoding/json"
	"fmt"
	"gosh/internal/program"
	"log"
	"os"
	"os/signal"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cache"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/websocket/v2"
)

type HttpServer struct {
	app                *fiber.App
	programHttpHandler *program.HttpHandler

	programSearchService *program.SearchService
}

func NewHttpServer(
	programHttpHandler *program.HttpHandler,
	programSearchService *program.SearchService,
) *HttpServer {
	app := fiber.New()

	return &HttpServer{
		app:                  app,
		programHttpHandler:   programHttpHandler,
		programSearchService: programSearchService,
	}
}

func (h *HttpServer) SetupHttpRoutes() {
	h.app.Post("/api/program", h.programHttpHandler.CreateProgramHandler)
	h.app.Get("/api/program", h.programHttpHandler.SearchProgram)
	h.app.Get("/api/program/:id", h.programHttpHandler.GetProgramDetailsHandler)
}

func (h *HttpServer) SetupWebSocketRoutes() {
	h.app.Get("/ws", websocket.New(func(c *websocket.Conn) {
		// c.Locals is added to the *websocket.Conn
		log.Println(c.Locals("allowed"))
		log.Println(c.Params("id"))
		log.Println(c.Query("v"))
		log.Println(c.Cookies("session"))

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

			programs := h.programSearchService.Search(context.Background(), req.Q)

			msg, err = json.Marshal(fiber.Map{
				"data": programs,
			})
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
}

func (h *HttpServer) SetupHttpMiddleware() {
	h.app.Use(logger.New())
	h.app.Use(cors.New())
}

func (h *HttpServer) SetupWebSocketMiddleware() {
	h.app.Use(cache.New(cache.Config{
		Next: func(c *fiber.Ctx) bool {
			return strings.Contains(c.Route().Path, "/ws")
		},
	}))

	h.app.Use("/ws", func(c *fiber.Ctx) error {
		// IsWebSocketUpgrade returns true if the client
		// requested upgrade to the WebSocket protocol.
		if websocket.IsWebSocketUpgrade(c) {
			c.Locals("allowed", true)
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})
}

func (h *HttpServer) RunWithGracefulShutdown() {
	c := make(chan os.Signal, 1)

	// os.Interrupt = Ctrl+C
	//  os.Kill = kill -9
	signal.Notify(c, os.Interrupt, os.Kill)

	serverShutdown := make(chan struct{})

	go func() {
		_ = <-c
		fmt.Println("Gracefully shutting down...")
		_ = h.app.Shutdown()
		serverShutdown <- struct{}{}
	}()

	if err := h.app.Listen(":3000"); err != nil {
		log.Panic(err)
	}

	<-serverShutdown

	fmt.Println("Running cleanup tasks...")
	fmt.Println("Done!")
}
