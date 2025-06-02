package router_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/sing3demons/todoapi/router"
	"github.com/sing3demons/todoapi/todo"
	"github.com/stretchr/testify/assert"
)

func TestFiberRouterGET(t *testing.T) {
	fiberApp := router.NewFiberRouter()

	fiberApp.GET("/ping", func(ctx todo.Context) {

		name := ctx.Param("name")

		ctx.JSON(http.StatusOK, DummyContext{Message: "pong" + name})
	})

	req := httptest.NewRequest(http.MethodGet, "/ping", nil)
	resp, _ := fiberApp.Test(req)

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var response DummyContext
	err := json.NewDecoder(resp.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Equal(t, "pong", response.Message)
}

func TestFiberRouterPOST(t *testing.T) {
	fiberApp := router.NewFiberRouter()

	fiberApp.POST("/echo", func(c todo.Context) {
		var payload DummyContext
		err := c.Bind(&payload)
		if err != nil {
			c.JSON(http.StatusBadRequest, DummyContext{Message: "invalid"})
			return
		}
		c.JSON(http.StatusOK, payload)
	})

	payload := DummyContext{Message: "hello"}
	body, _ := json.Marshal(payload)

	req := httptest.NewRequest(http.MethodPost, "/echo", strings.NewReader(string(body)))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := fiberApp.Test(req)

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var response DummyContext
	err := json.NewDecoder(resp.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Equal(t, payload.Message, response.Message)
}
