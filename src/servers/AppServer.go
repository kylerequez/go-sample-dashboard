package servers

import (
	"context"
	"fmt"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/static"
	"github.com/kylerequez/go-sample-dashboard/src/handlers"
)

type AppServer struct {
	Hostname string
	Port     string
}

type Server interface {
	Init()
}

func NewAppServer(hostname, port string) *AppServer {
	return &AppServer{
		Hostname: hostname,
		Port:     port,
	}
}

func (server *AppServer) Init() error {
	app := fiber.New(fiber.Config{
		AppName: "Dashboard App",
	})

	app.Use("/styles", static.New("./src/public/css"))
	app.Use("/htmx", static.New("./src/public/htmx"))

	ctx := context.TODO()
	if err := handlers.Init(app, ctx); err != nil {
		return err
	}

	return app.Listen(fmt.Sprintf("%s:%s", server.Hostname, server.Port))
}
