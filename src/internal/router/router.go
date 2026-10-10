package router

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/subhadip0x539/arc-dex-core-srv/src/internal/handler"
	"github.com/subhadip0x539/arc-dex-core-srv/src/internal/repository"
	"github.com/subhadip0x539/arc-dex-core-srv/src/internal/service"
)

type Router struct {
	pool *pgxpool.Pool
}

func NewRouter(pool *pgxpool.Pool) *Router {
	return &Router{pool: pool}
}

func (r *Router) Register(rg *gin.RouterGroup) {
	artifactRepository := repository.NewArtifactRepository(r.pool)
	artifactService := service.NewArtifactService(artifactRepository)
	artifactHandler := handler.NewArtifactHandler(artifactService)

	rg.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})

	v1 := rg.Group("/v1")

	artifact := v1.Group("/artifact")
	{
		artifact.POST("", artifactHandler.CreateArtifact)
	}

	artifacts := v1.Group("/artifacts")
	{
		artifacts.GET("", artifactHandler.GetArtifacts)
		artifacts.GET(":id", artifactHandler.GetArtifact)

		artifacts.POST(":id/domain", artifactHandler.CreateArtifact)
	}

}
