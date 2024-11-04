package server

import (
	"internal/handlers"
	"internal/repository"
	"internal/service"
	inPg "internal/storage/postgres"
	"log"
	"os"

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
	docHandler := handlers.MewDocumentHandler(dockService)

	// user  endpoints
	s.e.POST("/api/register", userHandlers.Register)
	s.e.POST("/api/auth", userHandlers.Authenticate)
	// s.e.DELETE() TODO: Доделать

	// document endpoints
	s.e.POST("/api/docs", docHandler.UploadDocument)
	// s.e.GET("/api/docs") TODO: Доделать
	// s.e.GET("/api/docs/:id") TODO: Доделать
	// s.e.DELETE("/api/docs/:id") TODO: Доделать

	// s.e.Use(middleware.AuthMiddleWare())
	s.e.Logger.Fatal(s.e.Start(":8080"))
}
