package infra

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type routeGroup struct {
	prefix  string
	handler func(rg *gin.RouterGroup)
}

type Server struct {
	engine *gin.Engine
	server *http.Server
	host   string
	port   int32
	cors   *cors.Config
	routes []routeGroup
}

func (s *Server) Host(host string) *Server {
	s.host = host
	return s
}

func (s *Server) Port(port int32) *Server {
	s.port = port
	return s
}

func (s *Server) CORS(cfg cors.Config) *Server {
	s.cors = &cfg
	return s
}

func (s *Server) RouteGroup(prefix string, fn func(*gin.RouterGroup)) *Server {
	s.routes = append(s.routes, routeGroup{prefix: prefix, handler: fn})
	return s
}

func (s *Server) Build() *Server {
	s.engine = gin.New()

	s.engine.Use(gin.Logger())
	s.engine.Use(gin.Recovery())

	if s.cors != nil {
		s.engine.Use(cors.New(*s.cors))
	} else {
		s.engine.Use(cors.New(cors.Config{
			AllowAllOrigins: true,
			AllowMethods:    []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
			AllowHeaders:    []string{"Origin", "Content-Type", "Authorization"},
			MaxAge:          12 * time.Hour,
		}))
	}

	for _, rg := range s.routes {
		rg.handler(s.engine.Group(rg.prefix))
	}

	s.server = &http.Server{
		Addr:    fmt.Sprintf("%s:%d", s.host, s.port),
		Handler: s.engine,
	}

	return s
}

func (s *Server) Start() error {
	if err := s.server.ListenAndServe(); err != nil {
		return fmt.Errorf("http server error: %w", err)
	}
	return nil
}

func (s *Server) Shutdown(ctx context.Context) error {
	if s.server == nil {
		return nil
	}

	return s.server.Shutdown(ctx)
}

func NewServer() *Server {
	return &Server{
		host: "0.0.0.0",
		port: 8080,
	}
}
