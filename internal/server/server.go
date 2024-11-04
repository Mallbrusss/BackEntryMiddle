package server

import (
	"log"
	"os"

	"github.com/Mallbrusss/BackEntryMiddle/internal/handlers"
	// customMiddleware "github.com/Mallbrusss/BackEntryMiddle/internal/middleware"
	"github.com/Mallbrusss/BackEntryMiddle/internal/repository"
	"github.com/Mallbrusss/BackEntryMiddle/internal/service"
	inPg "github.com/Mallbrusss/BackEntryMiddle/internal/storage/postgres"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
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

	//TODO: мб подредачить, но пока так
	uploadDir := "./../../uploads/documents"
	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		log.Fatalf("Ошибка при создании директории для загрузок: %v", err)
	}

	userRepo := repository.NewUserRepository(db)

	userService := service.NewUserService(userRepo)
	userHandlers := handlers.NewUserHandlers(userService)

	docRepo := repository.NewDocumentRepository(db)
	dockService := service.NewDocumentService(docRepo, uploadDir)
	docHandler := handlers.NewDocumentHandler(dockService)


	s.e.Use(middleware.Logger())
	s.e.Use(middleware.Recover())

	// user  endpoints
	s.e.POST("/api/register", userHandlers.Register)
	s.e.POST("/api/auth", userHandlers.Authenticate)
	// s.e.DELETE() TODO: Доделать

	// document endpoints
	s.e.POST("/api/docs", docHandler.UploadDocument)
	// s.e.GET("/api/docs") TODO: Доделать
	// s.e.GET("/api/docs/:id") TODO: Доделать
	// s.e.DELETE("/api/docs/:id") TODO: Доделать

	s.e.Logger.Fatal(s.e.Start(":8080"))
}
