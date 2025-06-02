package router_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/sing3demons/todoapi/router"
	"github.com/sing3demons/todoapi/todo"
	"github.com/stretchr/testify/assert"
)

type DummyContext struct {
	Message string `json:"message"`
}

func TestMyRouterGET(t *testing.T) {
	gin.SetMode(gin.TestMode)

	r := router.NewMyRouter()

	r.GET("/ping", func(ctx todo.Context) {

		name := ctx.Param("name")

		ctx.JSON(http.StatusOK, DummyContext{Message: "pong" + name})
	})

	req := httptest.NewRequest(http.MethodGet, "/ping", nil)
	w := httptest.NewRecorder()

	r.Test(w, req)

	res := w.Result()
	defer res.Body.Close()

	assert.Equal(t, http.StatusOK, res.StatusCode)

	var response DummyContext
	err := json.NewDecoder(res.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Equal(t, "pong", response.Message)
}

func TestMyRouterPOST(t *testing.T) {
	gin.SetMode(gin.TestMode)

	r := router.NewMyRouter()

	r.POST("/echo", func(c todo.Context) {
		var payload DummyContext
		err := c.Bind(&payload)
		if err != nil {
			c.JSON(http.StatusBadRequest, DummyContext{Message: "invalid"})
			return
		}
		c.JSON(http.StatusOK, payload)
	})

	body := `{"message":"hello"}`

	req := httptest.NewRequest(http.MethodPost, "/echo", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.Test(w, req)

	res := w.Result()
	defer res.Body.Close()

	assert.Equal(t, http.StatusOK, res.StatusCode)

	var response DummyContext
	err := json.NewDecoder(res.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Equal(t, "hello", response.Message)
}

func TestListenAndServeShutdown(t *testing.T) {
	os.Setenv("PORT", "3001")
	r := router.NewMyRouter()

	go func() {
		_ = r.ListenAndServe()
	}()

	// We don't connect, just test Shutdown
	err := r.Shutdown()
	assert.NoError(t, err)
}
