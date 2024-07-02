package handlers

import (
	"fmt"
	"strings"
	"time"

	"github.com/a-h/templ"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"github.com/gofiber/fiber/v2/middleware/cache"
)

// SetRoutes sets the routes in the application and
// associates them with their respective handlers
func SetRoutes(
	app *fiber.App, ah AuthHandler, sh SettingsHandler, sch SearchHandler,
) {
	// ↓ health checker route ↓
	app.Get("/health-checker", healthCheckerHandler)

	// ↓ Auth routes ↓
	app.Get("/login", ah.loginHandler)
	app.Post("/login", ah.loginPostHandler)

	// ↓ Admin routes ↓
	app.Get("/", ah.authMiddleware, sh.dashboardHandler)
	app.Post("/", ah.authMiddleware, sh.dashboardPostHandler)
	app.Post("/logout", ah.logoutHandler)

	// ↓ Create admin route [secret route] ↓
	app.Post("/create", ah.createAdminHandler)

	// ↓ Search route ↓
	app.Post("/search", sch.searchHandler)
	app.Use("/search", cache.New(cache.Config{
		Next: func(c *fiber.Ctx) bool {
			return c.Query("noCache") == "true"
		},
		Expiration:   30 * time.Minute,
		CacheControl: true,
	}))

	/* TODO: ↓ Fallback Page ↓ */
	app.Get("/*", func(c *fiber.Ctx) error {

		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Error 404: Not Found",
		})
	})
}

// healthCheckerHandler It is a handler that
// only checks that the server is online
func healthCheckerHandler(c *fiber.Ctx) error {
	name := strings.Title(c.Query("name"))
	if name == "" {
		name = "World"
	}

	return c.JSON(fiber.Map{
		"message": fmt.Sprintf("Hello, %s 😀!!", name),
	})
}

// Render creates the render of the
// functional template (go file)
// that we are passing to it
func Render(
	c *fiber.Ctx,
	component templ.Component,
	options ...func(*templ.ComponentHandler),
) error { // ↓ See NOTE-01 below ↓
	componentHandler := templ.Handler(component)
	for _, o := range options {
		o(componentHandler)
	}

	return adaptor.HTTPHandler(componentHandler)(c)
}

/* NOTE-01.- A-H/TEMPL INTEGRATION FIBER FRAMEWORK:
https://templ.guide/integrations/web-frameworks/#go-fiber
https://github.com/a-h/templ/blob/main/examples/integration-gofiber/main.go
*/

/* RENDERING HTML DIRECTLY:
https://templ.guide/static-rendering/blog-example#rendering-html-directly

func RenderHtmlString(html string) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) (err error) {
		_, err = io.WriteString(w, html)
		return
	})
}
*/
