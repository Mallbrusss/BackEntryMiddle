package server

import (
	"context"
	"log"
	"os"

	"github.com/Mallbrusss/BackEntryMiddle/internal/handlers"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	// customMiddleware "github.com/Mallbrusss/BackEntryMiddle/internal/middleware"
	"github.com/Mallbrusss/BackEntryMiddle/internal/repository"
	"github.com/Mallbrusss/BackEntryMiddle/internal/service"
	inPg "github.com/Mallbrusss/BackEntryMiddle/internal/storage/postgres"
	inRdb "github.com/Mallbrusss/BackEntryMiddle/internal/storage/redis"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Server struct {
	e         *echo.Echo
	db        *gorm.DB
	rdb       *redis.Client
	uploadDir string
}

func NewServer() *Server {
	return &Server{
		e:         echo.New(),
		uploadDir: "uploads/documents",
	}
}

func (s *Server) initializeDatabase() error {
	db, err := inPg.InitDB()
	if err != nil {
		return err
	}
	s.db = db
	return nil
}

func (s *Server) initializeRedis() error {
	rdb, err := inRdb.InitRedisCl()
	if err != nil {
		return err
	}
	s.rdb = rdb
	return nil
}

func (s *Server) createUploadDir() error {
	if err := os.MkdirAll(s.uploadDir, os.ModePerm); err != nil {
		return err
	}
	return nil
}

func (s *Server) initializeMiddleware() {
	s.e.Use(middleware.Logger())
	s.e.Use(middleware.Recover())
}

func (s *Server) registerRoutes() {
	// repository initialization
	userRepo := repository.NewUserRepository(s.db)
	docRepo := repository.NewDocumentRepository(s.db)

	// service initialization
	userService := service.NewUserService(userRepo)
	docService := service.NewDocumentService(docRepo, s.rdb, s.uploadDir)

	// handlers initialization
	userHandlers := handlers.NewUserHandlers(userService)
	docHandler := handlers.NewDocumentHandler(docService)

	// user  endpoints
	s.e.POST("/api/register", userHandlers.Register)
	s.e.POST("/api/auth", userHandlers.Authenticate)
	s.e.DELETE("/api/auth/:token", userHandlers.Logout)

	// document endpoints
	s.e.POST("/api/docs", docHandler.AuthMiddleWare()(docHandler.UploadDocument))
	s.e.GET("/api/docs", docHandler.AuthMiddleWare()(docHandler.GetDocuments))
	s.e.HEAD("/api/docs", docHandler.HeadDocument)
	s.e.GET("/api/docs/:id", docHandler.AuthMiddleWare()(docHandler.GetDocumentByID))
	s.e.HEAD("/api/docs/:id", docHandler.HeadDocument)
	s.e.DELETE("/api/docs/:id", docHandler.AuthMiddleWare()(docHandler.DeleteDocument))
}

func (s *Server) Shutdown(ctx context.Context) error {

	if s.rdb != nil {
		if err := s.rdb.Close(); err != nil {
			log.Printf("Redis close error: %v", err)
		}
	}
	return s.e.Shutdown(ctx)
}

func (s *Server) Run() error {
	if err := s.initializeDatabase(); err != nil {
		s.e.Logger.Fatal("Error connecting to database:", err)
	}

	if err := s.initializeRedis(); err != nil {
		s.e.Logger.Fatal("Error connecting to Redis:", err)
	}
	defer s.rdb.Close()

	if err := s.createUploadDir(); err != nil {
		s.e.Logger.Fatal("Error creating upload directory:", err)
	}

	s.initializeMiddleware()
	s.registerRoutes()

	return s.e.Start(":8080")
}
