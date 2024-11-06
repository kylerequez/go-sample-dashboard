package handlers

import (
	"github.com/gofiber/fiber/v3"
	"github.com/kylerequez/go-sample-dashboard/src/repositories"
	"github.com/kylerequez/go-sample-dashboard/src/utils"
	"github.com/kylerequez/go-sample-dashboard/src/views/pages"
)

type UserHandler struct {
	Ur  *repositories.UserRepository
	App *fiber.App
}

type userHandler interface {
	Init()
}

func NewUserHandler(
	app *fiber.App,
	ur *repositories.UserRepository,
) *UserHandler {
	return &UserHandler{
		App: app,
		Ur:  ur,
	}
}

func (uh *UserHandler) Init() {
	views := uh.App.Group("")
	views.Get("/users", uh.GetAllUsers)
}

func (uh *UserHandler) GetAllUsers(c fiber.Ctx) error {
	users, err := uh.Ur.GetAllUsers(c.Context())
	if err != nil {
		return err
	}

	return utils.Render(c, pages.Users(users))
}
