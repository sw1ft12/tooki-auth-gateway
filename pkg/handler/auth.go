package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"log"
	"net/http"
	"time"
	"tooki/pkg/authErrs"
	"tooki/pkg/models"
	"tooki/pkg/tokens"
)

func BindJSON(c *gin.Context, data any) *authErrs.Error {
	err := c.BindJSON(&data)
	if err != nil {
		return authErrs.New(authErrs.EINCORRECT, "некорректные данные"+err.Error(), "handler.BindJSON")
	}
	return nil
}

func Validate(v *validator.Validate, s any) *authErrs.Error {
	err := v.Struct(s)
	if err != nil {
		return authErrs.New(authErrs.EINCORRECT, err.Error(), "handler.Validate")
	}
	return nil
}

func SendError(c *gin.Context, err *authErrs.Error) {
	switch err.Code {
	case authErrs.EINTERNAL:
		c.JSON(http.StatusInternalServerError, err.Message)
	case authErrs.EEXIST:
		c.JSON(http.StatusConflict, err.Message)
	case authErrs.ENOTFOUND:
		c.JSON(http.StatusUnauthorized, err.Message)
	case authErrs.EINCORRECT:
		c.JSON(http.StatusBadRequest, err.Message)
	}
	logOut := fmt.Sprintf("\nType: %s\nMessage: %s\nOp: %s\n", err.Code, err.Message, err.Op)
	log.Println(logOut)
}

// @Summary		Регистрация пользователя
// @Description	Регистрация пользователя
// @Tags			Auth
// @Accept			json
// @Produce		json
// @Param			user	body		models.RegisterUserDto	true	"Данные для регистрации"
// @Success		201		{object}	models.RegisterResponse	"Пользователь зарегистрирован"
// @Failure		400		"Нверные данные"
// @Router			/register [post]
func (h *Handler) Register(c *gin.Context) {
	var dto models.RegisterUserDto
	err := BindJSON(c, &dto)
	if err != nil {
		SendError(c, err)
		return
	}

	err = Validate(h.validator, dto)
	if err != nil {
		SendError(c, err)
		return
	}

	user, err := h.repo.CreateUser(dto)
	if err != nil {
		SendError(c, err)
		return
	}

	c.JSON(http.StatusOK, user)
}

// @Summary		Аутентификация пользователя
// @Description	Аутентификация пользователя
// @Tags			Auth
// @Param			login	body	models.LoginUserDto	true	"Данные для аутентификации"
// @Accept			json
// @Produce		json
// @Success		200 {object} models.LoginResponse
// @Failure		400	"Неправильные логин или пароль"
// @Router			/login [post]
func (h *Handler) Login(c *gin.Context) {
	var dto models.LoginUserDto
	err := BindJSON(c, &dto)
	if err != nil {
		SendError(c, err)
		return
	}

	user, err := h.repo.GetUserByLogin(dto)
	if err != nil {
		SendError(c, err)
		return
	}

	refreshToken, accessToken, err := tokens.GenerateTokens(user)
	if err != nil {
		SendError(c, err)
		return
	}

	err = h.repo.SaveRefreshToken(refreshToken)
	if err != nil {
		SendError(c, err)
		return
	}

	c.SetCookie("refresh_token", refreshToken.Token, int(time.Now().Add(time.Hour*24*30).Unix()), "/", "", true, true)
	resp := models.LoginResponse{
		AccessToken: accessToken,
		Id:          user.Id,
		Name:        user.Name,
		Email:       user.Email,
		Gender:      user.Gender,
		Role:        user.Role,
	}
	c.JSON(http.StatusOK, resp)
}

// @Summary		Выход из учётной записи
// @Description	Выход из учётной записи и удаление refresh_token пользователя
// @Tags			Auth
// @Success		200 "Сброс refresh_token пользователя"
// @Router			/logout [get]
func (h *Handler) LogOut(c *gin.Context) {
	c.SetCookie("refresh_token", "", -1, "/", "localhost", false, true)
}
