package router

import (
	"log"
	"net/http"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	_ "github.com/gofiber/fiber/v2/middleware/monitor"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/sing3demons/todoapi/todo"
)

type FiberCtx struct {
	*fiber.Ctx
}

func NewFiberCtx(c *fiber.Ctx) *FiberCtx {
	return &FiberCtx{Ctx: c}
}

func (c *FiberCtx) Bind(out interface{}) error {
	return c.Ctx.BodyParser(out)
}

func (c *FiberCtx) JSON(code int, v interface{}) {
	err := c.Ctx.Status(code).JSON(v)
	if err != nil {
		log.Println(err)
		return
	}
}

func (c *FiberCtx) Param(key string) string {
	return c.Ctx.Params(key)
}
func (c *FiberCtx) Status(code int) {
	c.Ctx.Status(code)
}

//Router

type FiberRouter struct {
	*fiber.App
}

func NewFiberRouter() *FiberRouter {
	app := fiber.New()
	app.Use(recover.New())

	// app.Get("/dashboard", monitor.New())
	app.Use(logger.New(logger.ConfigDefault))
	return &FiberRouter{app}
}

func NewFiberHandler(handler func(todo.Context)) fiber.Handler {
	return func(c *fiber.Ctx) error {
		handler(NewFiberCtx(c))
		return nil
	}
}

func (r *FiberRouter) POST(path string, handler func(todo.Context)) {
	r.App.Post(path, NewFiberHandler(handler))
}

func (r *FiberRouter) GET(path string, handler func(todo.Context)) {
	r.App.Get(path, NewFiberHandler(handler))
}

func (r *FiberRouter) ListenAndServe() error {
	return r.App.Listen(":" + os.Getenv("PORT"))
}

func (r *FiberRouter) Shutdown() error {
	return r.App.Shutdown()
}

func (r *FiberRouter) Test(req *http.Request, msTimeout ...int) (resp *http.Response, err error) {
	return r.App.Test(req, msTimeout...)
}
