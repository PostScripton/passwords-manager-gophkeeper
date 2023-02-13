package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/PostScripton/passwords-manager-gophkeeper/internal/models"
	"github.com/PostScripton/passwords-manager-gophkeeper/internal/repository"
	"golang.org/x/sync/errgroup"
	"math/rand"
	"time"
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

	rand.Seed(time.Now().UnixMicro())

	query := `INSERT INTO creds_secrets (uid, website, login, enc_password, additional_data, user_id) VALUES ($1, $2, $3, $4, $5, $6);`
	if _, err = r.db.ExecContext(ctx, query, rand.Int63(), website, login, encPassword, additionalData, userID); err != nil {
		return fmt.Errorf("store creds secret to the SQLite database on create: %w", err)
	}

	return nil
}

func (r *CredsSecretsRepository) GetById(ctx context.Context, uid int64) (*models.CredsSecret, error) {
	query := `SELECT id, uid, website, login, enc_password, additional_data, user_id FROM creds_secrets WHERE uid = $1;`

	creds := new(models.CredsSecret)
	if err := r.db.QueryRowContext(ctx, query, uid).Scan(
		&creds.ID,
		&creds.UID,
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

func (r *CredsSecretsRepository) Delete(ctx context.Context, uid int64) error {
	query := `DELETE FROM creds_secrets WHERE uid = $1;`

	if _, err := r.db.ExecContext(ctx, query, uid); err != nil {
		return fmt.Errorf("deleting creds from SQLite: %w", err)
	}

	return nil
}

func (r *CredsSecretsRepository) GetList(ctx context.Context, userID int) ([]*models.CredsSecret, error) {
	query := `SELECT id, uid, website, login, enc_password, additional_data, user_id
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
			&secret.UID,
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
		// [Спринт 10: Примитивы синхронизации: Пакеты sync и x/sync]
		// Делаем специально так, чтобы в горутину не попали последние значения цикла
		uid := secret.UID

		deleteGroup.Go(func() error {
			return r.Delete(deleteCtx, uid)
		})
	}
	if err := deleteGroup.Wait(); err != nil {
		return err
	}

	for _, secret := range list {
		query := `INSERT INTO creds_secrets (uid, website, login, enc_password, additional_data, user_id) VALUES ($1, $2, $3, $4, $5, $6);`
		if _, err := r.db.ExecContext(
			ctx,
			query,
			secret.UID,
			secret.Website,
			secret.Login,
			secret.Password,
			secret.AdditionalData,
			secret.UserID,
		); err != nil {
			return fmt.Errorf("store creds secret to the SQLite database on set list: %w", err)
		}
	}

	return nil
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
