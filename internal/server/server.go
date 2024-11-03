package server

import (
	"internal/repository"
	"internal/service"
	"internal/handlers"
	inPg "internal/storage/postgres"

	"github.com/labstack/echo/v4"
)

type Server struct {
	e *echo.Echo
}

func NewServer() *Server {
	return &Server{
		e: echo.New(),
	}
}

func (s *Server) Run() {

	db := inPg.InitDB()
	userRepo := repository.NewUserRepository(db)

	userService := service.NewUserService(userRepo)
	userHandlers := handlers.NewUserHandlers(userService)

	s.e.POST("/api/register", userHandlers.Register)
	s.e.POST("/api/auth", userHandlers.Authenticate)

	// s.e.Use(middleware.AuthMiddleWare())
	s.e.Logger.Fatal(s.e.Start(":8080"))
}
