package repository

import (
	"context"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
	"tooki/pkg/authErrs"
	"tooki/pkg/models"
)

type Repo struct {
	pool *pgxpool.Pool
}

func NewRepo(pool *pgxpool.Pool) *Repo {
	return &Repo{pool: pool}
}

func (r *Repo) CreateUser(data models.RegisterUserDto) (models.RegisterResponse, *authErrs.Error) {
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

	op := "repository.CreateUser"
	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			return models.RegisterResponse{}, authErrs.New(authErrs.EINTERNAL, pgErr.Message, op)
		}
		return models.RegisterResponse{}, authErrs.New(authErrs.EEXIST, "Пользователь уже существует", op)
	}

	query := `SELECT id, email, name, role, created_at, verified, banned FROM Users`
	row, _ := r.pool.Query(ctx, query)
	resp, _ := pgx.CollectOneRow(row, pgx.RowToStructByNameLax[models.RegisterResponse])
	return resp, nil
}

func (r *Repo) GetUserByLogin(dto models.LoginUserDto) (*models.User, *authErrs.Error) {
	ctx := context.Background()
	row, err := r.pool.Query(ctx, `SELECT * FROM Users WHERE login = @login`,
		pgx.NamedArgs{
			"login": dto.Login,
		})

	op := "repository.GetUserByLogin"
	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			return nil, authErrs.New(authErrs.EINTERNAL, pgErr.Message, op)
		}
		return nil, authErrs.New(authErrs.EEXIST, "Неправильный логин или пароль", op)
	}

	user, _ := pgx.CollectOneRow(row, pgx.RowToStructByNameLax[models.User])

	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(dto.Password)) != nil {
		return nil, authErrs.New(authErrs.EEXIST, "Неправильный логин или пароль", op)
	}

	return &user, nil
}

func (r *Repo) SaveRefreshToken(token *models.RefreshToken) *authErrs.Error {
	ctx := context.Background()
	query := `INSERT INTO Tokens (user_id, token, expires_in) VALUES (@user_id, @token, @expires_in) 
			ON CONFLICT(user_id) DO UPDATE SET token=@token, expires_in=@expires_in`
	_, err := r.pool.Query(ctx, query, pgx.NamedArgs{
		"userId":     token.UserId,
		"token":      token.Token,
		"expires_in": token.ExpiresIn,
	})

	op := "repository.Tokens"
	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			return authErrs.New(authErrs.EINTERNAL, pgErr.Message, op)
		}
		return authErrs.New(authErrs.EINTERNAL, "Невалидный токен", op)
	}
	return nil
}
