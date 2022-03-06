package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"github.com/sing3demons/todoapi/database"
	"github.com/sing3demons/todoapi/router"
	"github.com/sing3demons/todoapi/store"
	"github.com/sing3demons/todoapi/todo"
)

func init() {
	err := godotenv.Load("dev.env")
	if err != nil {
		log.Printf("please consider environment variables: %s", err)
	}
}

var (
	buildcommit = "dev"
	buildtime   = time.Now().String()
)

func main() {
	_, err := os.Create("/tmp/live")
	if err != nil {
		log.Fatal(err)
	}
	defer os.Remove("/tmp/live")

	r := router.NewMyRouter()

	r.GET("/", func(ctx todo.Context) {
		ctx.JSON(200, map[string]interface{}{
			"message": "hello world",
		})
	})

	r.GET("/healthz", func(ctx todo.Context) {
		ctx.Status(200)
	})

	r.GET("x", func(ctx todo.Context) {
		ctx.JSON(200, map[string]interface{}{
			"buildcommit": buildcommit,
			"buildtime":   buildtime,
		})
	})

	database.InitDB()
	db := database.GetDB()
	store := store.NewGormStore(db)

	// collection := database.Collection()
	// store := store.NewMongoStore(collection)
	handler := todo.NewTodoHandler(store)

	r.POST("/todos", handler.NewTask)
	r.GET("/todos", handler.List)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	serve := &http.Server{
		Addr:           ":" + os.Getenv("PORT"),
		Handler:        r,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	go func() {
		if err := serve.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	<-ctx.Done()
	stop()
	fmt.Println("shutting down gracefully, press Ctrl+C again to force")

	timeoutCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := serve.Shutdown(timeoutCtx); err != nil {
		fmt.Println(err)
	}
}
