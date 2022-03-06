package router

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/sing3demons/todoapi/todo"
)

type MyContext struct {
	*gin.Context
}

func NewMyContext(ctx *gin.Context) *MyContext {
	return &MyContext{Context: ctx}
}

func (c *MyContext) Bind(v interface{}) error {
	return c.Context.ShouldBindJSON(v)
}
func (c *MyContext) JSON(code int, v interface{}) {
	c.Context.JSON(code, v)
}

func (c *MyContext) Param(key string) string {
	return c.Context.Param(key)
}

func NewMyHandler(handler func(todo.Context)) gin.HandlerFunc {
	return func(c *gin.Context) {
		handler(NewMyContext(c))
	}
}

type MyRouter struct {
	*gin.Engine
}

func NewMyRouter() *MyRouter {
	r := gin.Default()
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{
		"http://localhost:8081",
	}
	config.AllowHeaders = []string{
		"Origin",
		"Authorization",
	}

	r.Use(cors.New(config))
	return &MyRouter{r}
}

// POST(relativePath string, handlers ...HandlerFunc)
func (r *MyRouter) POST(path string, handler func(todo.Context)) {
	r.Engine.POST(path, NewMyHandler(handler))
}

func (r *MyRouter) GET(path string, handler func(todo.Context)) {
	r.Engine.GET(path, NewMyHandler(handler))
}
