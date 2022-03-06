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
	// r := router.NewFiberRouter()

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

	//Graceful Shutdown
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go func() {
		if err := r.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	<-ctx.Done()
	stop()
	fmt.Println("shutting down gracefully, press Ctrl+C again to force")

	if err := r.Shutdown(); err != nil {
		fmt.Println(err)
	}
}
