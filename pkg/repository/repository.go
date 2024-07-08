package repository

import (
	"context"
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
	"time"
	"tooki/pkg/models"
)

type Repo struct {
	pool *pgxpool.Pool
}

func NewRepo(pool *pgxpool.Pool) *Repo {
	return &Repo{pool: pool}
}

type RegisterUserDto struct {
	Email    string `json:"email"  binding:"required" example:"test@email.com"`
	Login    string `json:"login" binding:"required" example:"sw1ft12"`
	Password string `json:"password" binding:"required" example:"assag23214"`
	Name     string `json:"name" binding:"required" example:"Артёмчик Zиновьев"`
	Age      int    `json:"age" example:"5"`
	Gender   string `json:"gender"  example:"Female"`
}

type RegisterResponse struct {
	Id        string    `json:"id" db:"id" example:"8cbabbe9-5fff-4dbe-a77e-104bf4e63dbe"`
	Email     string    `json:"email" db:"email" example:"test@gmail.com"`
	Name      string    `json:"name" db:"name" example:"Зиновьев Артём"`
	Role      string    `json:"role" db:"role" example:"USER"`
	CreatedAt time.Time `json:"created_at" db:"created_at" example:"2024-03-02"`
	Verified  bool      `json:"verified" db:"verified"`
	Banned    bool      `json:"banned" db:"banned"`
}

func (r *Repo) CreateUser(data RegisterUserDto) (RegisterResponse, error) {

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(data.Password), 12)

	ctx := context.Background()
	registerQuery := `INSERT INTO Users (email, login, password, name, age, gender) VALUES 
                        (@email, @login, @password, @name, @age, @gender)`
	_, err := r.pool.Query(ctx, registerQuery, pgx.NamedArgs{
		"email":    data.Email,
		"login":    data.Login,
		"password": hashedPassword,
		"name":     data.Name,
		"age":      data.Age,
		"gender":   data.Gender,
	})

	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			return RegisterResponse{}, pgErr
		}
		return RegisterResponse{}, errors.New("такой пользователь уже существует")
	}

	query := `SELECT id, email, name, role, created_at, verified, banned FROM Users`
	row, err := r.pool.Query(ctx, query)
	resp, _ := pgx.CollectOneRow(row, pgx.RowToStructByNameLax[RegisterResponse])
	return resp, nil
}

type LoginUserDto struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type LoginResponse struct {
	AccessToken string `json:"token" db:"token"`
	Id          string `json:"id" db:"id"`
	Name        string `json:"name" db:"name"`
	Email       string `json:"email"`
	Gender      string `json:"gender"`
	Role        string `json:"role"`
}

func (r *Repo) GetUserByLogin(data LoginUserDto) (*models.User, error) {
	ctx := context.Background()
	row, err := r.pool.Query(ctx, `SELECT * FROM Users WHERE login = @login`,
		pgx.NamedArgs{"login": data.Login})
	if err != nil {
		return nil, errors.New("неправильный логин или пароль")
	}

	user, _ := pgx.CollectOneRow(row, pgx.RowToStructByNameLax[models.User])

	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(data.Password)) != nil {
		return nil, errors.New("неправильный логин или пароль")
	}

	return &user, nil
}

type RefreshToken struct {
	Token          string
	ExpirationTime time.Time
}

func (r *Repo) CreateRefreshToken(userId string) error {
	ctx := context.Background()
	query := "INSERT INTO Tokens (user_id) VALUES (@userId)"
	_, err := r.pool.Query(ctx, query, pgx.NamedArgs{
		"userId": userId,
	})
	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			return pgErr
		}
	}
	return nil
}

func (r *Repo) UpdateRefreshToken(token *jwt.Token) {
	q := `UPDATE Tokens SET refresh_token = @token`
	ctx := context.Background()
	r.pool.QueryRow(ctx, q, pgx.NamedArgs{
		"token": token.Raw,
	})
}
