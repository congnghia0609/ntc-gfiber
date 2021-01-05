/**
 *
 * @author nghiatc
 * @since Jan 5, 2021
 */

package server

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/congnghia0609/ntc-gconf/nconf"
	"github.com/go-playground/validator"
	"github.com/gofiber/fiber"
	"github.com/gofiber/fiber/middleware/limiter"
	"github.com/gofiber/fiber/middleware/pprof"
	"github.com/gofiber/fiber/middleware/recover"
	"github.com/gofiber/fiber/middleware/timeout"
	"github.com/gofiber/template/html"
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

// Timeout handler
func Timeout(c *fiber.Ctx) error {
	time.Sleep(5 * time.Second)
	return c.SendString("This is message from timeout handler!")
}

// Home handler
func Home(c *fiber.Ctx) error {
	// Render index template
	return c.Render("layout", fiber.Map{
		"Title": "Home Page",
	})
}

//// Validation:
// https://github.com/go-playground/validator

// Job struct
type Job struct {
	Type   string `json:"type" validate:"required,min=3,max=32"`
	Salary int    `json:"salary" validate:"required,number"`
}

// User struct
type User struct {
	Name     string `json:"name" validate:"required,min=3,max=32"`
	IsActive bool   `json:"isactive" validate:"required,eq=True|eq=False"`
	Email    string `json:"email" validate:"required,email,min=6,max=32"`
	Job      Job    `json:"job" validate:"dive"`
}

// ErrorResponse struct
type ErrorResponse struct {
	FailedField string
	Tag         string
	Value       string
}

// ValidateStruct function
func ValidateStruct(user User) []*ErrorResponse {
	var errors []*ErrorResponse
	validate := validator.New()
	err := validate.Struct(user)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			var element ErrorResponse
			element.FailedField = err.StructNamespace()
			element.Tag = err.Tag()
			element.Value = err.Param()
			errors = append(errors, &element)
		}
	}
	return errors
}

// AddUser function
func AddUser(c *fiber.Ctx) error {
	//Connect to database
	user := new(User)
	if err := c.BodyParser(user); err != nil {
		c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
		return err
	}
	ju, _ := json.Marshal(user)
	log.Println("user:", string(ju))
	errs := ValidateStruct(*user)
	if errs != nil {
		return c.JSON(errs)
		// mes, _ := json.Marshal(errs)
		// return errors.New(string(mes))
	}
	//Do something else here

	//Return user
	return c.JSON(user)
}

// StartWebServer start WebServer
// https://github.com/gofiber/fiber
func StartWebServer(name string) {
	// Config
	c := nconf.GetConfig()
	address := c.GetString(name + ".addr")

	// 1. Setup Fiber
	// app := fiber.New()
	// Custom config Fiber: https://docs.gofiber.io/api/fiber
	// Initialize standard Go html template engine: https://github.com/gofiber/template
	// engine: html | ace | amber | django | handlebars | jet | mustache | pug
	engine := html.New("./views", ".html")
	fcfg := fiber.Config{
		Prefork:         false,           // Enables use of the SO_REUSEPORT socket option. This will spawn multiple Go processes listening on the same port. Default: false
		ServerHeader:    "NTC",           // Enables the Server HTTP header with the given value. Default: ""
		BodyLimit:       4 * 1024 * 1024, // Sets the maximum allowed size for a request body. Default: 4 * 1024 * 1024
		Concurrency:     256 * 1024,      // Maximum number of concurrent connections. Default: 256 * 1024
		ReadBufferSize:  4096,            // Per-connection buffer size for requests' reading. This also limits the maximum header size. Default: 4096
		WriteBufferSize: 4096,            // Per-connection buffer size for responses' writing. Default: 4096
		Views:           engine,          // Template Middleware for supported engines. Default: nil
	}
	app := fiber.New(fcfg)

	// 2. Setup Module Middleware
	// 2.1. Limiter used to limit repeated requests to public APIs
	// app.Use(limiter.New())
	// Custom config
	app.Use(limiter.New(limiter.Config{
		// Next: func(c *fiber.Ctx) bool {
		// 	// While list IP => skip Limiter
		// 	log.Println("Limiter IP:", c.IP())
		// 	return c.IP() == "127.0.0.1"
		// }, // Default: nill
		Max:        600,             // Max connections / Expiration (second). Default: 5 Reqs/Expiration
		Expiration: 1 * time.Minute, // Expiration is the time to count requests in memory. Default: 1 * time.Minute
	}))

	// 2.2. Recover used to recovers from panics anywhere
	// https://docs.gofiber.io/api/middleware/recover
	app.Use(recover.New())

	// 2.3. Pprof visualization tool
	// http://localhost:8080/debug/pprof/
	app.Use(pprof.New())

	// 2.4. Timeout middleware set timeout for handler.
	// https://docs.gofiber.io/api/middleware/timeout
	// http://localhost:8080/timeout
	app.Get("/timeout", timeout.New(Timeout, 3*time.Second))

	// 3. Setup Router
	// https://docs.gofiber.io/api/app
	// http://localhost:8080/hello/nghia
	app.Get("/hello/:name", Hello)
	// http://localhost:8080/json/nghia
	app.Get("/json/:name", RespJSON)
	// http://localhost:8080/panic
	app.Get("/panic", func(c *fiber.Ctx) error {
		// This panic will be catch by the middleware Recover.
		panic("I'm an error")
	})

	// 3.1. Group routes
	// https://docs.gofiber.io/guide/grouping
	api := app.Group("/api", Auth) // /api
	v1 := api.Group("/v1")         // /api/v1
	v1.Get("/list/:name", Hello)   // /api/v1/list/nghia
	v1.Get("/user/:name", Hello)   // /api/v1/user/nghia
	v2 := api.Group("/v2")         // /api/v2
	v2.Get("/list/:name", Hello)   // /api/v2/list/nghia
	v2.Get("/user/:name", Hello)   // /api/v2/user/nghia

	// 3.2. Template engines: https://github.com/gofiber/template
	// http://localhost:8080/home
	app.Get("/home", Home)

	// 3.3. Validation: https://docs.gofiber.io/guide/validation
	// Running a test with the following curl commands
	// curl -X POST -H "Content-Type: application/json" --data "{\"name\":\"john\",\"isactive\":true}" http://localhost:8080/register/user
	app.Post("/register/user", AddUser)

	// 4. Static Files Handler
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
