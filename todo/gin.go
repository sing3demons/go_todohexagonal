package todo

import "github.com/gin-gonic/gin"

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

func NewMyHandler(handler func(Context)) gin.HandlerFunc {
	return func(c *gin.Context) {
		handler(NewMyContext(c))
	}
}
