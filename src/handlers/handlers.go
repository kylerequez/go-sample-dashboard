package handlers

import (
	"context"

	"github.com/gofiber/fiber/v3"
	"github.com/kylerequez/go-sample-dashboard/src/db"
	"github.com/kylerequez/go-sample-dashboard/src/repositories"
)

func Init(app *fiber.App, ctx context.Context) error {
	if err := db.Connect(ctx); err != nil {
		return err
	}

	if err := db.Ping(ctx); err != nil {
		return err
	}

	ur := repositories.NewUserRepository(db.DB, "users")

	uh := NewUserHandler(app, ur)
	uh.Init()

	ah := NewAuthHandler(app, ur)
	ah.Init()

	return nil
}
