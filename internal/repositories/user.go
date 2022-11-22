package repositories

import (
	"context"
	"fmt"

	"gitlab.privy.id/go_graphql/internal/entity"
	"gitlab.privy.id/go_graphql/pkg/postgres"
)

type User interface {
	RegisterUser(context.Context, *entity.User) error
}

type userImplementation struct {
	conn postgres.Adapter
}

// Create User
func (r userImplementation) RegisterUser(ctx context.Context, user *entity.User) error {
	query :=
		`
			INSERT INTO users (

				name,
				email,
				password,
				created_at,
				deleted_at
			)
			VALUES ($1, $2, $3, $4, $5) RETURNING id
		`

	err := r.conn.QueryRow(
		ctx,
		query,
		user.Name,
		user.Email,
		user.Password,
		user.CreatedAt,
		user.DeletedAt,
	).Scan(
		&user.ID,
	)

	if err != nil {
		fmt.Println(err)
	}
	return nil
}

func NewUserImplementation(conn postgres.Adapter) *userImplementation {
	return &userImplementation{
		conn: conn,
	}
}
