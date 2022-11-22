package repositories

import (
	"context"
	"fmt"

	"gitlab.privy.id/go_graphql/internal/entity"
	"gitlab.privy.id/go_graphql/pkg/postgres"
)

type Shot interface {
	Get(context.Context) ([]entity.Shot, error)
	GetShotID(context.Context, int) (entity.Shot, error)
	CreateShot(context.Context, *entity.Shot) (*entity.Shot, error)
	DeleteShot(context.Context, int) error
	UpdateShot(context.Context, entity.Shot) error
}

type shotImplementation struct {
	conn postgres.Adapter
}

// Get ALl Shot
func (r shotImplementation) Get(ctx context.Context) ([]entity.Shot, error) {
	query := `
		SELECT
			id,
			title,
			img,
			description,
			category
		FROM shots
	`
	rows, err := r.conn.QueryRows(ctx, query)
	if err != nil {
		return nil, err
	}

	shots := []entity.Shot{}

	for rows.Next() {
		var shot entity.Shot

		err = rows.Scan(
			&shot.ID,
			&shot.Title,
			&shot.IMG,
			&shot.Description,
			&shot.Category,
		)
		if err != nil {
			return nil, err
		}
		shots = append(shots, shot)
	}

	// fmt.Println(shots)

	return shots, nil
}

// Get Shot By ID
func (r shotImplementation) GetShotID(ctx context.Context, id int) (shot entity.Shot, err error) {
	query := `
		SELECT
			id,
			title,
			img,
			description,
			category
		FROM shots WHERE id=$1
	`
	row := r.conn.QueryRow(ctx, query, id)
	err = row.Scan(
		&shot.ID,
		&shot.Title,
		&shot.IMG,
		&shot.Description,
		&shot.Category,
	)
	if err != nil {
		err = fmt.Errorf("Scanning a user: %w", err)
		return entity.Shot{}, err
	}

	return shot, nil
}

// Create Shot
func (r shotImplementation) CreateShot(ctx context.Context, shot *entity.Shot) (*entity.Shot, error) {
	query :=

		`INSERT INTO shots (
		title,
		img,
		description,
		category,
		created_at,
		updated_at,
		deleted_at
	)
	VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING *
	`

	var output entity.Shot
	err := r.conn.QueryRow(
		ctx,
		query,
		shot.Title,
		shot.IMG,
		shot.Description,
		shot.Category,
		shot.CreatedAt,
		shot.UpdatedAt,
		shot.DeletedAt,
	).Scan(
		&output.ID,
		&output.Title,
		&output.IMG,
		&output.Description,
		&output.Category,
		&output.CreatedAt,
		&output.UpdatedAt,
		&output.DeletedAt,
	)

	if err != nil {
		err = fmt.Errorf("executing query: %w", err)
		return nil, err
	}

	return &output, nil
}

// Delete Shot
func (r shotImplementation) DeleteShot(ctx context.Context, id int) error {
	query := `
		DELETE FROM shots
		WHERE id = $1
	`

	_, err := r.conn.Exec(ctx, query, id)
	if err != nil {
		err = fmt.Errorf("Deleting shot: %w", err)

		return err
	}
	return nil
}

// Update Shot
func (r shotImplementation) UpdateShot(ctx context.Context, shot entity.Shot) error {
	query := `
		UPDATE shots SET
			title = $1,
			img = $2,
			description = $3,
			category = $4,
			updated_at = $5
			WHERE id = $6
	`

	_, err := r.conn.Exec(
		ctx,
		query,
		&shot.Title,
		&shot.IMG,
		&shot.Description,
		&shot.Category,
		&shot.UpdatedAt,
		&shot.ID,
	)

	if err != nil {
		err = fmt.Errorf("updating shot: %w", err)
		return err
	}
	return nil
}

func NewShotImplementation(conn postgres.Adapter) *shotImplementation {
	return &shotImplementation{
		conn: conn,
	}
}
