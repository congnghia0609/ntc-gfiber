/**
 *
 * @author nghiatc
 * @since Jan 5, 2021
 */

package server

import (
	"fmt"
	"log"
	"time"

	"github.com/congnghia0609/ntc-gconf/nconf"
	"github.com/gofiber/fiber"
	"github.com/gofiber/fiber/middleware/pprof"
)

// Hello handler
// https://docs.gofiber.io/api/ctx
func Hello(c *fiber.Ctx) error {
	if !fiber.IsChild() {
		fmt.Println("I'm the parent process")
	} else {
		fmt.Println("I'm a child process")
	}
	msg := fmt.Sprintf("Hello, %s!", c.Params("name"))
	if c.Locals("role") == "admin" {
		msg = fmt.Sprintf("Hello, %s! (admin)", c.Params("name"))
	}
	return c.SendString(msg)
}

// RespJSON handler
// https://docs.gofiber.io/api/ctx#json
func RespJSON(c *fiber.Ctx) error {
	if !fiber.IsChild() {
		fmt.Println("I'm the parent process")
	} else {
		fmt.Println("I'm a child process")
	}
	mapData := fiber.Map{
		"name": c.Params("name"),
		"time": time.Now(),
	}
	return c.JSON(mapData)
}

// Auth authentication handler
func Auth(c *fiber.Ctx) error {
	log.Println("In Auth Handler...")
	c.Locals("role", "admin")
	return c.Next()
}

// StartWebServer start WebServer
// https://github.com/gofiber/fiber
func StartWebServer(name string) {
	// Config
	c := nconf.GetConfig()
	address := c.GetString(name + ".addr")

	// Setup Fiber
	// app := fiber.New()
	// Custom config: https://docs.gofiber.io/api/fiber
	fcfg := fiber.Config{
		Prefork:         false,           // Enables use of the SO_REUSEPORT socket option. This will spawn multiple Go processes listening on the same port. Default: false
		ServerHeader:    "NTC",           // Enables the Server HTTP header with the given value. Default: ""
		BodyLimit:       4 * 1024 * 1024, // Sets the maximum allowed size for a request body. Default: 4 * 1024 * 1024
		Concurrency:     256 * 1024,      // Maximum number of concurrent connections. Default: 256 * 1024
		ReadBufferSize:  4096,            // Per-connection buffer size for requests' reading. This also limits the maximum header size. Default: 4096
		WriteBufferSize: 4096,            // Per-connection buffer size for responses' writing. Default: 4096
	}
	app := fiber.New(fcfg)

	// Setup Module Middleware
	// 1. Pprof visualization tool
	// http://localhost:8080/debug/pprof/
	app.Use(pprof.New())

	// Setup Router
	// https://docs.gofiber.io/api/app
	// http://localhost:8080/hello/nghia
	app.Get("/hello/:name", Hello)
	// http://localhost:8080/json/nghia
	app.Get("/json/:name", RespJSON)

	// Group routes
	// https://docs.gofiber.io/api/app#group
	api := app.Group("/api", Auth) // /api
	v1 := api.Group("/v1")         // /api/v1
	v1.Get("/list/:name", Hello)   // /api/v1/list/nghia
	v1.Get("/user/:name", Hello)   // /api/v1/user/nghia
	v2 := api.Group("/v2")         // /api/v2
	v2.Get("/list/:name", Hello)   // /api/v2/list/nghia
	v2.Get("/user/:name", Hello)   // /api/v2/user/nghia

	// Static Files Handler
	// https://docs.gofiber.io/api/app#static
	app.Static("/static", "./public")
	// => http://localhost:8080/static/css/main.css
	// app.Static("/", "./public/dist")
	// => http://localhost:8080/index.html
	// => http://localhost:8080/js/script.js
	// => http://localhost:8080/css/style.css
	// app.Static("*", "./public/dist/index.html")
	// => http://localhost:8080/any/path/shows/index/html

	log.Printf("======= WebServer[%s] is running on host: %s", name, address)
	log.Fatal(app.Listen(address))
}
