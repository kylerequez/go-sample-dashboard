package handlers

import (
	"errors"
	"net/mail"

	"github.com/gofiber/fiber/v3"
	"github.com/kylerequez/go-sample-dashboard/src/models"
	"github.com/kylerequez/go-sample-dashboard/src/repositories"
	"github.com/kylerequez/go-sample-dashboard/src/utils"
	"github.com/kylerequez/go-sample-dashboard/src/views/pages"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	App *fiber.App
	Ur  *repositories.UserRepository
}

type authHandler interface {
	Init()
}

type LoginForm struct {
	Email    string
	Password string
}

type SignupForm struct {
	Name     string
	Email    string
	Password string
}

func NewAuthHandler(
	app *fiber.App,
	ur *repositories.UserRepository,
) *AuthHandler {
	return &AuthHandler{
		App: app,
		Ur:  ur,
	}
}

func (ah *AuthHandler) Init() {
	views := ah.App.Group("")
	views.Get("/", ah.GetHomePage)
	views.Get("/login", ah.GetLoginPage)
	views.Get("/sign-up", ah.GetSignupPage)

	api := ah.App.Group("/api/v1/auth")
	api.Post("/login", ah.LoginUser)
	api.Post("/sign-up", ah.SignupUser)
}

func (ah *AuthHandler) GetHomePage(c fiber.Ctx) error {
	info := models.AppInfo{
		Title:       "HOME",
		CurrentPath: c.Path(),
	}

	return utils.Render(c, pages.Home(info))
}

func (ah *AuthHandler) GetLoginPage(c fiber.Ctx) error {
	info := models.AppInfo{
		Title:       "LOGIN",
		CurrentPath: c.Path(),
	}

	form := models.LoginFormData{
		Errors: make(map[string]error),
	}

	return utils.Render(c, pages.Login(info, form))
}

func (ah *AuthHandler) LoginUser(c fiber.Ctx) error {
	form := models.LoginFormData{
		Errors: make(map[string]error),
	}

	body := new(LoginForm)
	if err := c.Bind().Body(body); err != nil {
		form.Errors["FORM"] = errors.New("unable to retrieve form input")
		return utils.Render(c, pages.LoginForm(form))
	}

	form.Email = body.Email
	form.Password = body.Password

	hasError := false
	if body.Email == "" {
		hasError = true
		form.Errors["EMAIL"] = errors.New("email is empty")
	}

	_, err := mail.ParseAddress(body.Email)
	if err != nil {
		hasError = true
		form.Errors["EMAIL"] = errors.New("email is invalid")
	}

	if body.Password == "" {
		hasError = true
		form.Errors["PASSWORD"] = errors.New("password is empty")
	}

	if hasError {
		return utils.Render(c, pages.LoginForm(form))
	}

	user, err := ah.Ur.GetUserByEmail(c.Context(), body.Email)
	if err != nil {
		form.Errors["FORM"] = err
		return utils.Render(c, pages.LoginForm(form))
	}

	if err := bcrypt.CompareHashAndPassword(user.Password, []byte(body.Password)); err != nil {
		form.Errors["PASSWORD"] = err
		return utils.Render(c, pages.LoginForm(form))
	}

	c.Set("HX-Redirect", "true")
	return c.Redirect().To("/admin/dashboard")
}

func (ah *AuthHandler) GetSignupPage(c fiber.Ctx) error {
	info := models.AppInfo{
		Title:       "SIGN UP",
		CurrentPath: c.Path(),
	}

	form := models.SignupFormData{
		Errors: make(map[string]error),
	}

	return utils.Render(c, pages.Signup(info, form))
}

func (ah *AuthHandler) SignupUser(c fiber.Ctx) error {
	form := models.SignupFormData{
		Errors: make(map[string]error),
	}

	body := new(SignupForm)
	if err := c.Bind().Body(body); err != nil {
		form.Errors["FORM"] = err
		return utils.Render(c, pages.SignupForm(form))
	}

	form.Name = body.Name
	form.Email = body.Email
	form.Password = body.Password

	hasError := false
	if body.Name == "" {
		hasError = true
		form.Errors["NAME"] = errors.New("name is empty")
	}

	if body.Email == "" {
		hasError = true
		form.Errors["EMAIL"] = errors.New("email is empty")
	}

	_, err := mail.ParseAddress(body.Email)
	if err != nil {
		hasError = true
		form.Errors["EMAIL"] = errors.New("email is invalid")
	}

	if body.Password == "" {
		hasError = true
		form.Errors["PASSWORD"] = errors.New("password is empty")
	}

	if hasError {
		return utils.Render(c, pages.SignupForm(form))
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(body.Password), bcrypt.DefaultCost)
	if err != nil {
		form.Errors["PASSWORD"] = err
		return utils.Render(c, pages.SignupForm(form))
	}

	user := models.User{
		Name:     body.Name,
		Email:    body.Email,
		Password: hashedPassword,
	}
	if err := ah.Ur.CreateUser(c.Context(), user); err != nil {
		form.Errors["FORM"] = err
		return utils.Render(c, pages.SignupForm(form))
	}

	return utils.Render(c, pages.SignupForm(form))
}
