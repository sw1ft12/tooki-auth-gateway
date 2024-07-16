package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5/pgxpool"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"regexp"
	"tooki/pkg/repository"
)

type Handler struct {
	repo      *repository.Repo
	validator *validator.Validate
}

func NewValidator() *validator.Validate {
	v := validator.New(validator.WithRequiredStructEnabled())
	_ = v.RegisterValidation("password", func(fl validator.FieldLevel) bool {
		ok, _ := regexp.MatchString(`(?=.*[a-z])(?=.*[A-Z])(?=.*\\d)(?=.*[@$!%*?&])[A-Za-z\\d@$!%*?&]{8,}`, fl.Field().String())
		return ok
	})
	return v
}

func NewHandler(pool *pgxpool.Pool) *Handler {
	return &Handler{
		repo:      repository.NewRepo(pool),
		validator: NewValidator(),
	}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.Default()
	router.POST("/login", h.Login)
	router.POST("/register", h.Register)
	router.GET("/logout", h.LogOut)
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	return router
}
