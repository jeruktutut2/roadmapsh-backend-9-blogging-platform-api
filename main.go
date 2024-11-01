package main

import (
	"blogging-platform-api/controllers"
	"blogging-platform-api/repositories"
	"blogging-platform-api/routes"
	"blogging-platform-api/services"
	"blogging-platform-api/utils"
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

func main() {
	postgresUtil := utils.NewPostgresConnection()
	validate := validator.New()
	e := echo.New()

	blogRepository := repositories.NewBlogRepository()
	blogService := services.NewBlogService(postgresUtil, validate, blogRepository)
	blogController := controllers.NewBlogController(blogService)
	routes.BlogRoute(e, blogController)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()
	go func() {
		if err := e.Start(os.Getenv("ECHO_HOST")); err != nil && err != http.ErrServerClosed {
			e.Logger.Fatal("shutting down the server")
		}
	}()

	<-ctx.Done()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
}
