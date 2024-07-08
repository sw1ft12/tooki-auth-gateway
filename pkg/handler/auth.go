package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"log"
	"net/http"
	"time"
	"tooki/pkg/models"
	"tooki/pkg/repository"
)

// @Summary		Регистрация пользователя
// @Description	Регистрация пользователя
// @Tags			Auth
// @Accept			json
// @Produce		json
// @Param			user	body		repository.RegisterUserDto	true	"Данные для регистрации"
// @Success		200		{object}	repository.RegisterResponse	"Пользователь зарегистрирован"
// @Failure		400		"Нверные данные"
// @Router			/register [post]
func (h *Handler) Register(c *gin.Context) {
	var data repository.RegisterUserDto
	err := c.BindJSON(&data)
	if err != nil {
		c.JSON(http.StatusBadRequest, "некорректные данные"+err.Error())
		return
	}

	user, err := h.repo.CreateUser(data)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	c.JSON(http.StatusOK, user)
}

type tokenClaims struct {
	jwt.RegisteredClaims
	Id     string `json:"id"`
	Name   string `json:"name"`
	Email  string `json:"email"`
	Gender string `json:"gender"`
	Role   string `json:"role"`
}

const key = "22323"

func GenerateTokens(user *models.User) (repository.RefreshToken, string, error) {
	accessToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 5)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
		Id:     user.Id,
		Name:   user.Name,
		Email:  user.Email,
		Gender: user.Gender,
		Role:   user.Role,
	}).SignedString([]byte(key))
	refreshToken := repository.RefreshToken{Token: uuid.NewString(), ExpirationTime: time.Now().Add(time.Hour * 30 * 24)}
	return refreshToken, accessToken, err
}

// @Summary		Аутентификация пользователя
// @Description	Аутентификация пользователя
// @Tags			Auth
// @Param			login	body	repository.LoginUserDto	true	"Данные для аутентификации"
// @Accept			json
// @Produce		json
// @Success		200
// @Failure		400	"Неправильные логин или пароль"
// @Router			/login [post]
func (h *Handler) Login(c *gin.Context) {
	var data repository.LoginUserDto
	err := c.BindJSON(&data)
	if err != nil {
		c.JSON(http.StatusBadRequest, "Некорректные данные")
		return
	}
	user, err := h.repo.GetUserByLogin(data)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	refreshToken, accessToken, err := GenerateTokens(user)
	if err != nil {
		log.Fatal(err)
	}

	err = h.repo.CreateRefreshToken(user.Id)
	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			c.JSON(http.StatusInternalServerError, pgErr.Error())
			return
		}
	}

	c.SetCookie("refresh_token", refreshToken.Token, int(time.Now().Add(time.Hour*24*30).Unix()), "/", "", true, true)
	resp := repository.LoginResponse{
		AccessToken: accessToken,
		Id:          user.Id,
		Name:        user.Name,
		Email:       user.Email,
		Gender:      user.Gender,
		Role:        user.Role,
	}
	c.JSON(http.StatusOK, resp)
}
