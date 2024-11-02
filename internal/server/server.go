package server

import (
	"net/http"

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


func (s *Server) Run(){
	s.e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "hello test")
	})

	s.e.Logger.Fatal(s.e.Start(":8080"))
}