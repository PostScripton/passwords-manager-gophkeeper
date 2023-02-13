package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/PostScripton/passwords-manager-gophkeeper/internal/models"
	"github.com/PostScripton/passwords-manager-gophkeeper/internal/repository"
	"golang.org/x/sync/errgroup"
)

type CredsSecretsRepository struct {
	db *SQLite
}

var _ repository.CredsSecrets = (*CredsSecretsRepository)(nil)

func NewCredsSecretsRepository(db *SQLite) repository.CredsSecrets {
	return &CredsSecretsRepository{
		db: db,
	}
}

func (r *CredsSecretsRepository) Create(
	ctx context.Context,
	userID int,
	website, login, encPassword, additionalData string,
) error {
	exists, err := r.checkCredsSecretExists(ctx, userID, website, login)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("credentials for this website exist")
	}

	query := `INSERT INTO creds_secrets (website, login, enc_password, additional_data, user_id) VALUES ($1, $2, $3, $4, $5);`
	if _, err = r.db.ExecContext(ctx, query, website, login, encPassword, additionalData, userID); err != nil {
		return fmt.Errorf("store creds secret to the database: %w", err)
	}

	return nil
}

func (r *CredsSecretsRepository) GetById(ctx context.Context, id int64) (*models.CredsSecret, error) {
	query := `SELECT id, website, login, enc_password, additional_data, user_id FROM creds_secrets WHERE id = $1;`

	creds := new(models.CredsSecret)
	if err := r.db.QueryRowContext(ctx, query, id).Scan(
		&creds.ID,
		&creds.Website,
		&creds.Login,
		&creds.Password,
		&creds.AdditionalData,
		&creds.UserID,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, err
	}

	return creds, nil
}

func (r *CredsSecretsRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM creds_secrets WHERE id = $1;`

	if _, err := r.db.ExecContext(ctx, query, id); err != nil {
		return fmt.Errorf("deleting creds from SQLite: %w", err)
	}

	return nil
}

func (r *CredsSecretsRepository) GetList(ctx context.Context, userID int) ([]*models.CredsSecret, error) {
	query := `SELECT id, website, login, enc_password, additional_data, user_id
		FROM creds_secrets
		WHERE user_id = $1
		ORDER BY website, login;`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}

	list := make([]*models.CredsSecret, 0)

	for rows.Next() {
		secret := new(models.CredsSecret)
		if err = rows.Scan(
			&secret.ID,
			&secret.Website,
			&secret.Login,
			&secret.Password,
			&secret.AdditionalData,
			&secret.UserID,
		); err != nil {
			return nil, err
		}

		list = append(list, secret)
	}

	return list, nil
}

func (r *CredsSecretsRepository) SetList(ctx context.Context, list []models.CredsSecret) error {
	deleteGroup, deleteCtx := errgroup.WithContext(ctx)
	for _, secret := range list {
		deleteGroup.Go(func() error {
			return r.Delete(deleteCtx, secret.ID)
		})
	}
	if err := deleteGroup.Wait(); err != nil {
		return err
	}

	createGroup, createCtx := errgroup.WithContext(ctx)
	for _, secret := range list {
		createGroup.Go(func() error {
			exists, err := r.checkCredsSecretExists(createCtx, secret.UserID, secret.Website, secret.Login)
			if err != nil {
				return err
			}
			if exists {
				return fmt.Errorf("credentials for this website exist")
			}

			// Для того чтобы не автоинкрементил ID и не получались разные секреты
			// (13 - 1... + 1 (auto-increment) = 13)
			id := secret.ID - 1

			query := `INSERT INTO creds_secrets (id, website, login, enc_password, additional_data, user_id) VALUES ($1, $2, $3, $4, $5, $6);`
			if _, err = r.db.ExecContext(
				createCtx,
				query,
				id,
				secret.Website,
				secret.Login,
				secret.Password,
				secret.AdditionalData,
				secret.UserID,
			); err != nil {
				return fmt.Errorf("store creds secret to the database: %w", err)
			}

			return nil
		})
	}

	return createGroup.Wait()
}

func (r *CredsSecretsRepository) Truncate(ctx context.Context) error {
	query := `DELETE FROM creds_secrets;`

	if _, err := r.db.ExecContext(ctx, query); err != nil {
		return err
	}

	return nil
}

func (r *CredsSecretsRepository) checkCredsSecretExists(
	ctx context.Context,
	userID int,
	website, login string,
) (bool, error) {
	query := `SELECT COUNT(*) FROM creds_secrets WHERE website = $1 and login = $2 and user_id = $3;`

	var count int
	if err := r.db.QueryRowContext(ctx, query, website, login, userID).Scan(&count); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}

		return false, err
	}

	return count > 0, nil
}
