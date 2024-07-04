package handlers

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/emarifer/search-engine/internal/handlers/dto"
	"github.com/emarifer/search-engine/internal/services"
	"github.com/emarifer/search-engine/internal/utils"
	"github.com/emarifer/search-engine/views"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

/********** Handlers for Auth Views **********/

type AuthService interface {
	CreateAdmin(ad services.User) error
	LoginAsAdmin(email, password string) (services.User, error)
}

func NewAuthHandler(us AuthService) AuthHandler {

	return AuthHandler{
		UserServices: us,
	}
}

type AuthHandler struct {
	UserServices AuthService
}

func (ah *AuthHandler) loginHandler(c *fiber.Ctx) error {

	return Render(c, views.Login())
}

func (ah *AuthHandler) loginPostHandler(c *fiber.Ctx) error {
	userLogin := dto.LoginFormDto{
		Email:     strings.Trim(c.FormValue("email"), " "),
		Passsword: strings.Trim(c.FormValue("password"), " "),
	}

	// time.Sleep(2 * time.Second) // to see the loading indicator
	// Check email & password
	if userLogin.Email == "" || userLogin.Passsword == "" {

		return c.
			Status(fiber.StatusBadRequest).
			SendString("✖&nbsp;&nbsp; email or password cannot be empty")
	}

	// Find credentials as admin in DB
	user, err := ah.UserServices.LoginAsAdmin(
		userLogin.Email, userLogin.Passsword,
	)
	if err != nil {

		return c.
			Status(fiber.StatusUnauthorized).
			SendString("✖&nbsp;&nbsp; unauthorized user")
	}

	// fmt.Printf("%+v\n", user)

	// Create JWT
	signedToken, err := utils.CreateNewAuthToken(
		user.ID, user.Email, user.IsAdmin,
	)
	if err != nil {

		return c.
			Status(fiber.StatusInternalServerError).
			SendString("✖&nbsp;&nbsp; something went wrong logging in, please try again")
	}

	// Create and set the cookie
	cookie := fiber.Cookie{
		Name:     "admin",
		Value:    signedToken,
		Expires:  time.Now().Add(24 * time.Hour),
		HTTPOnly: true, // Meant only for the server
	}
	c.Cookie(&cookie)

	// c.Status(fiber.StatusSeeOther).Set("HX-Location", "/") // 303 code

	// return nil
	return c.Status(fiber.StatusSeeOther).Redirect("/")
}

func (ah *AuthHandler) logoutHandler(c *fiber.Ctx) error {
	c.ClearCookie("admin")

	return c.Status(fiber.StatusSeeOther).Redirect("/login")
}

func (ah *AuthHandler) createAdminHandler(c *fiber.Ctx) error {
	ad := services.User{
		Email:    "your@email.com",
		Password: "your-password",
		IsAdmin:  true,
	}
	err := ah.UserServices.CreateAdmin(ad)
	if err != nil {

		return c.
			Status(fiber.StatusInternalServerError).
			SendString(fmt.Sprintf("something went wrong: %s", err))
	}

	return c.
		Status(fiber.StatusCreated).
		SendString("Admin Created!!")
}

/********** Handlers for Dashboard Views **********/

type SettingsService interface {
	Get() (services.SearchSettings, error)
	Upadate(amount uint, searchOn, addNew bool) error
}

func NewSettingsHandler(ss SettingsService) SettingsHandler {

	return SettingsHandler{
		SearchConfig: ss,
	}
}

type SettingsHandler struct {
	SearchConfig SettingsService
}

func (sh *SettingsHandler) dashboardHandler(c *fiber.Ctx) error {

	settings, err := sh.SearchConfig.Get()
	if err != nil {

		return c.
			Status(fiber.StatusInternalServerError).
			SendString("✖&nbsp;&nbsp; something went wrong")
	}

	amount := strconv.FormatUint(uint64(settings.Amount), 10)

	return Render(c, views.Home(amount, settings.SearchOn, settings.AddNew))
}

func (sh *SettingsHandler) dashboardPostHandler(c *fiber.Ctx) error {
	// time.Sleep(2 * time.Second) // to see the loading indicator

	settings := dto.SettingsFormDto{}
	if err := c.BodyParser(&settings); err != nil {

		return c.
			Status(fiber.StatusInternalServerError).
			SendString("✖&nbsp;&nbsp; something went wrong")
	}

	if settings.Amount == 0 {

		return c.
			Status(fiber.StatusBadRequest).
			SendString("✖&nbsp;&nbsp; amount cannot be empty")
	}

	err := sh.SearchConfig.Upadate(
		settings.Amount, settings.SearchOn, settings.AddNew,
	)
	if err != nil {

		return c.
			Status(fiber.StatusInternalServerError).
			SendString("✖&nbsp;&nbsp; something went wrong")
	}

	return c.Status(fiber.StatusSeeOther).Redirect("/")
}

func (ah *AuthHandler) authMiddleware(c *fiber.Ctx) error {

	// Get the cookie by name
	cookie := c.Cookies("admin")
	if cookie == "" {

		return c.Status(fiber.StatusSeeOther).Redirect("/login")
	}

	// Parse the cookie & check for errors
	token, err := jwt.ParseWithClaims(
		cookie,
		&utils.AuthClaims{},
		func(t *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("SECRET_KEY")), nil
		},
	)
	if err != nil {

		return c.Status(fiber.StatusSeeOther).Redirect("/login")
	}

	// Parse the custom claims & check jwt is valid
	_, ok := token.Claims.(*utils.AuthClaims)
	if ok && token.Valid {
		return c.Next()
	}

	return c.Status(fiber.StatusSeeOther).Redirect("/login")
}

/*
Printing Struct Variables in Golang:
https://www.geeksforgeeks.org/printing-struct-variables-in-golang/

*/
