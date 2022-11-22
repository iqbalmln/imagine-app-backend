package repositories

import (
	"context"
	"fmt"

	"gitlab.privy.id/go_graphql/internal/entity"
	"gitlab.privy.id/go_graphql/pkg/logger"
	"gitlab.privy.id/go_graphql/pkg/postgres"
	"golang.org/x/crypto/bcrypt"
)

type Auth interface {
	Get(context.Context) ([]entity.Login, error)
	Register(context.Context, *entity.Login) error
	HashPassword(password string) (string, error)
	CheckPasswordHash(password, hash string) bool
	Login(context.Context, string) (entity.Login, error)
}

type loginImplement struct {
	conn postgres.Adapter
}

func (l loginImplement) Login(ctx context.Context, email string) (entity.Login, error) {
	var user entity.Login
	query := `SELECT 
			id,
			username,
			email,
			password,
			created_at
			FROM users WHERE email = $1
	`

	err := l.conn.QueryRow(ctx, query, email).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Password,
		// &user.Role,
		&user.CreatedAt,
	)
	if err != nil {
		logger.Error(err)
		return entity.Login{}, err
	}

	return user, nil
}

func (l loginImplement) Register(ctx context.Context, login *entity.Login) error {
	query := `
		INSERT INTO users(
			username,
			email,
			password,
			created_at
	)
	VALUES ($1,$2,$3,$4) returning id
	`
	err := l.conn.QueryRow(
		ctx,
		query,
		login.Username,
		login.Email,
		login.Password,
		// login.Role,
		login.CreatedAt).Scan(&login.ID)

	if err != nil {
		err = fmt.Errorf("executing query: %w", err)
		return err
	}
	return nil
}

func (l loginImplement) Get(ctx context.Context) ([]entity.Login, error) {
	query := `
		SELECT
			*
		FROM users
	`

	// fmt.Println(query)
	rows, err := l.conn.QueryRows(ctx, query)

	if err != nil {
		fmt.Println("Error querying")
	}

	users := []entity.Login{}

	for rows.Next() {
		var user entity.Login

		err = rows.Scan(
			&user.ID,
			&user.Username,
			&user.Email,
			&user.Password,
			&user.CreatedAt,
		)

		// fmt.Println(user)
		if err != nil {
			fmt.Println("gagal scan")
			// return nil, err
		}

		users = append(users, user)
	}

	return users, nil
}

func (l loginImplement) HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func (l loginImplement) CheckPasswordHash(hash, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func NewLoginImplment(conn postgres.Adapter) *loginImplement {
	return &loginImplement{
		conn: conn,
	}
}
