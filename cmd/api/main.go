package main

import (
	"log"
	"os"

	"my-todo/internal/adapter/handler"
	"my-todo/internal/adapter/router"
	"my-todo/internal/infra/mysql"
	"my-todo/internal/usecase"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("no .env file found")
	}

	dsn := os.Getenv("DATABASE_URL")
	client, err := mysql.NewClient(dsn)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer client.Close()

	repo := mysql.NewTodoRepository(client)
	uc := usecase.NewTodoUsecase(repo)
	h := handler.NewTodoHandler(uc)

	r := router.SetupRouter(h)
	r.Run(":8080")
}
