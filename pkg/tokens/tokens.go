package tokens

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"os"
	"time"
	"tooki/pkg/authErrs"
	"tooki/pkg/models"
)

type tokenClaims struct {
	jwt.RegisteredClaims
	Id     string `json:"id"`
	Name   string `json:"name"`
	Email  string `json:"email"`
	Gender string `json:"gender"`
	Role   string `json:"role"`
}

func GenerateTokens(user *models.User) (*models.RefreshToken, string, *authErrs.Error) {
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
	}).SignedString([]byte(os.Getenv("JWT_SECRET")))
	refreshToken := models.RefreshToken{
		UserId:    user.Id,
		Token:     uuid.NewString(),
		ExpiresIn: time.Now().Add(time.Hour * 30 * 24),
	}
	if err != nil {
		return &refreshToken, accessToken, authErrs.New(authErrs.EINTERNAL, err.Error(), "tokens.GenerateTokens")
	}
	return &refreshToken, accessToken, nil
}
