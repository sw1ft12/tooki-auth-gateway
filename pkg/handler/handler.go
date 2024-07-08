package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"tooki/pkg/repository"
)

type Handler struct {
	repo *repository.Repo
}

func NewHandler(pool *pgxpool.Pool) *Handler {
	return &Handler{repo: repository.NewRepo(pool)}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.Default()

	router.POST("/login", h.Login)
	router.POST("/register", h.Register)

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	return router
}
