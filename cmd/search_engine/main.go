package main

import (
	"fmt"
	"log"
	"os"

	// "os/signal"
	// "syscall"
	"time"

	"github.com/emarifer/search-engine/db"
	"github.com/emarifer/search-engine/internal/handlers"
	"github.com/emarifer/search-engine/internal/services"
	"github.com/emarifer/search-engine/internal/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
)

func main() {
	env := godotenv.Load()
	if env != nil {
		panic("cannot find environment variables from file")
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = ":5500"
	} else {
		port = fmt.Sprintf(":%s", port)
	}

	app := fiber.New(fiber.Config{
		IdleTimeout: 5 * time.Second, // ↓ See note below ↓
	})
	app.Static("/", "./assets")
	app.Use(logger.New())
	app.Use(compress.New())

	// Init database
	db.InitDB()

	// Dependency injection
	as := services.NewAdminServices(services.User{}, db.GetDB())
	ah := handlers.NewAuthHandler(&as)
	ss := services.NewSearchSettingsServices(
		services.SearchSettings{}, db.GetDB(),
	)
	sh := handlers.NewSettingsHandler(&ss)
	si := services.NewSearchIndexServices(services.SearchIndex{}, db.GetDB())
	sch := handlers.NewSearchHandler(&si)

	handlers.SetRoutes(app, ah, sh, sch)

	us := services.NewUrlServices(services.CrawledUrl{}, db.GetDB())
	is := services.NewSearchIndexServices(services.SearchIndex{}, db.GetDB())
	utils.StartCronJobs(&ss, &us, &is)

	log.Println("🚀 Starting server and listening at port", port)

	log.Fatal(app.Listen(port))
}

// 🧬 ⚡️ 🚀 🏁

/* REFERENCES:
https://youtu.be/9ziQpDWH70I?si=IAq2NOWSkwP5vNmB&t=239

The response-targets Extension:
https://v1.htmx.org/extensions/response-targets/

BUILD A WEB CRAWLER IN GO:
https://juliensalinas.com/en/how-to-speed-up-web-scraping-with-go-golang-concurrency/
https://golangbot.com/buffered-channels-worker-pools/
https://jackdanger.com/build-a-web-crawler-in-go.html
https://www.google.com/search?q=golang+web+crawler+concurrency&sca_esv=a5295d27d6bcdf79&sca_upv=1&sxsrf=ADLYWILsjN7Io1rO6sXmIXAkSRdrNOFqUw:1719917840873&ei=EN2DZuP1NLeI9u8P3NiewAc&start=0&sa=N&sstk=Ad9T53ybVqhkuanCOn8RxYXZAUFn-_H9RwpaTQ2a5zkk47z9P8vKXSwZNiEjnlawYOE8_76npoZmg0GwOMN1NBpsqv4kKImPk3e_TMIHcvy1a8VM3jn_3YWai_xQwohU5XWE&ved=2ahUKEwjjvcqfmYiHAxU3hP0HHVysB3g4ChDy0wN6BAgCEAQ&biw=1280&bih=919&dpr=1

GOLANG NAMING RULES AND CONVENTIONS:
https://freedium.cfd/https://medium.com/@kdnotes/golang-naming-rules-and-conventions-8efeecd23b68

AKHIL SHARMA REPOSITORY AND VIDEO:
https://github.com/AkhilSharma90/go-full-text-search
https://www.youtube.com/watch?v=BPLpzpgp79A

How to Pretty Print JSON output with cURL:
https://mkyong.com/web/how-to-pretty-print-json-output-in-curl/

curl -v -X POST http://localhost:5500/search -d '{ "term": "partner" }' -H "content-type: application/json" | python3 -m json.tool

*/

/* APPLICATION SHUTDOWN CONFIGURATION ALTERNATIVE:
// Start our server and listen for a shutdown
go func() {
	if err := app.Listen(port); err != nil {
		log.Panic(err)
	}
}()

c := make(chan os.Signal, 1)
signal.Notify(c, os.Interrupt, syscall.SIGTERM)

<-c // Block the main thread until interupted
app.Shutdown()
log.Println("Shutting down server")
*/

/* `GOPLS` SERVER CRASH WITH `D-E/TEMPLE` EXTENSION IN VSCODE:

The extension does not work, at least on Linux, with version v0.16.0.
You need to downgrade to v0.14.0.

Manual download:
https://stackoverflow.com/questions/76867499/how-to-set-the-version-of-gopls-to-install-with-vscode-go-extension
https://github.com/golang/tools/tree/master/gopls

*/
