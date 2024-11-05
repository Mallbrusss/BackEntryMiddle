package server

import (
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
	e   *echo.Echo
	db  *gorm.DB
	rdb *redis.Client
}

func NewServer() *Server {
	return &Server{
		e:   echo.New(),
		db:  inPg.InitDB(),
		rdb: inRdb.InitRedisCl(),
	}
}

func (s *Server) Run() {

	// create dir
	uploadDir := "uploads/documents"
	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		log.Fatalf("Error create download dir: %v", err)
	}

	// service initialization
	userRepo := repository.NewUserRepository(s.db)
	userService := service.NewUserService(userRepo)
	userHandlers := handlers.NewUserHandlers(userService)
	docRepo := repository.NewDocumentRepository(s.db)
	docService := service.NewDocumentService(docRepo, s.rdb, uploadDir)
	docHandler := handlers.NewDocumentHandler(docService)

	// using middleware
	// s.e.Use(middleware.Logger())
	s.e.Use(middleware.Recover())

	// user  endpoints
	s.e.POST("/api/register", userHandlers.Register)
	s.e.POST("/api/auth", userHandlers.Authenticate)
	s.e.DELETE("/api/auth/:token", userHandlers.Logout)

	// document endpoints
	s.e.POST("/api/docs", docHandler.AuthMiddleWare()(docHandler.UploadDocument))
	s.e.GET("/api/docs", docHandler.AuthMiddleWare()(docHandler.GetDocuments))
	// s.e.HEAD("/api/docs", )// TODO: Доделать
	s.e.GET("/api/docs/:id", docHandler.AuthMiddleWare()(docHandler.GetDocumentByID))
	// s.e.HEAD("/api/docs/:id", ) // TODO: Доделать
	s.e.DELETE("/api/docs/:id", docHandler.AuthMiddleWare()(docHandler.DeleteDocument))

	s.e.Logger.Fatal(s.e.Start(":8080"))
}
