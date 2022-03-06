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
	"github.com/sing3demons/todoapi/router"
	"github.com/sing3demons/todoapi/store"
	"github.com/sing3demons/todoapi/todo"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
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

	db, err := gorm.Open(mysql.Open(os.Getenv("DB_CONN")), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	db.AutoMigrate(&todo.Todo{})

	uri := "mongodb://root:admin1234@localhost:27017"
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}
	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		panic(err)
	}

	collection := client.Database("go-todo").Collection("todos")
	


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

	gormStore := store.NewGormStore(db)
	_ = gormStore
	mongo := store.NewMongoStore(collection)
	handler := todo.NewTodoHandler(mongo)

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
